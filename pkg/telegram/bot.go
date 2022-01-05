package telegram

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"qbot"
	"qbot/pkg/service"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	bot        *tgbotapi.BotAPI
	service    *service.Service
	goCron     *gocron.Scheduler
	adminsList []int64
}

func NewBot(bot *tgbotapi.BotAPI, service *service.Service, goCron *gocron.Scheduler, adminsList []int64) *Bot {
	return &Bot{
		bot:        bot,
		service:    service,
		goCron:     goCron,
		adminsList: adminsList,
	}
}

// Start запуск бота
func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	var updates tgbotapi.UpdatesChannel
	webhookHost := os.Getenv("WEBHOOK_HOST")

	if webhookHost == "" {
		// Запуск бота в режиме long polling
		_, err := b.bot.RemoveWebhook()
		if err != nil {
			return err
		}
		updates, err = b.bot.GetUpdatesChan(u)
		if err != nil {
			return err
		}
		log.Println("Bot started on long polling...")
	} else {
		// Запуск бота на вебхуак, если указана переменная WEBHOOK_HOST
		_, err := b.bot.RemoveWebhook()
		if err != nil {
			return err
		}
		log.Printf("Setting webhook on %s...\n", webhookHost)
		_, err = b.bot.SetWebhook(tgbotapi.NewWebhook(webhookHost + b.bot.Token))
		if err != nil {
			log.Fatalf("Setting webhook error: %s", err.Error())
			return err
		}
		log.Println("Getting webhook info...")
		info, err := b.bot.GetWebhookInfo()
		if err != nil {
			return err
		}
		if info.LastErrorDate != 0 {
			return fmt.Errorf("telegram callback failed: %s", info.LastErrorMessage)
		}
		log.Println("Getting updates channel...")
		updates = b.bot.ListenForWebhook("/" + b.bot.Token)
		go http.ListenAndServe("localhost:8012", nil)
		log.Println("Bot started on webhook...")
	}

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			b.handleError(0, errors.New("unknown behaviour"))
			continue
		}
		if update.CallbackQuery != nil {
			if err := b.handleQuery(update.CallbackQuery); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}
		if update.Message.IsCommand() {
			log.Printf("Take command\n")
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}
		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
	return nil
}

// minNumber найти минимальное значение
func minNumber(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

// runGoroutinesForMassMailing запуск горутин для массовой рассылки
func (b *Bot) runGoroutinesForMassMailing(content *[]qbot.Answer, messagesChan chan tgbotapi.Message) {
	var wg sync.WaitGroup
	maxMessagePerSecond := 30
	maxMessageSendAsync := make(chan struct{}, 20)
	for i, elem := range *content {
		wg.Add(1)
		_elem := elem
		go func(messagesChan chan tgbotapi.Message, wg *sync.WaitGroup) {
			maxMessageSendAsync <- struct{}{}
			for {
				sliceIndex := minNumber(len(_elem.Content), 50)
				log.Printf("Send message to %d: (%s)...", _elem.ChatId, _elem.Content[:sliceIndex])
				message, err := b.SendMessage(_elem)
				if err == nil {
					<-maxMessageSendAsync
				}
				if err != nil {
					matched, _ := regexp.MatchString("Too Many Requests: retry after", err.Error())
					if matched {
						re := regexp.MustCompile("[0-9]+")
						seconds, _ := strconv.Atoi(re.FindAllString(err.Error(), -1)[0])
						log.Printf("%s. Sleep", err.Error())
						time.Sleep(time.Duration(seconds + 1))
						log.Printf("Error send message to %d \"%s\". Retry send...", _elem.ChatId, err.Error())
						continue
					}
					log.Printf("Error: send message to %d %s", _elem.ChatId, err.Error())
					messagesChan <- tgbotapi.Message{MessageID: 0, Chat: &tgbotapi.Chat{ID: _elem.ChatId}}
					wg.Done()
					break
				}
				log.Printf("Message to %d sended", _elem.ChatId)
				wg.Done()
				messagesChan <- message
				return
			}
		}(messagesChan, &wg)
		if i % maxMessagePerSecond == 0 {
			time.Sleep(time.Second * 2)
		}
	}
	wg.Wait()
}

// MassMailing функция для массовой рассылки
func (b *Bot) MassMailing(content []qbot.Answer) ([]int64, error) {
	mailingId, err := b.service.CreateMailing()
	contentLength := len(content)
	if err != nil {
		return []int64{}, err
	}
	log.Printf("Start mailing (%d)...", mailingId)
	messagesChan := make(chan tgbotapi.Message, contentLength)
	b.runGoroutinesForMassMailing(&content, messagesChan)

	var sendedMessages []tgbotapi.Message
	var sendedToChatIds []int64
	var deactivatedSubscriptionIds []int64
	var counter int
	for c := range messagesChan {
		if c.MessageID != 0 {
			sendedMessages = append(sendedMessages, c)
			sendedToChatIds = append(sendedToChatIds, c.Chat.ID)
		} else {
			deactivatedSubscriptionIds = append(deactivatedSubscriptionIds, c.Chat.ID)
		}
		counter++
		if counter == contentLength {
			break
		}
	}

	log.Printf("Mailing (%d) sended", mailingId)
	if len(deactivatedSubscriptionIds) > 0 {
		if err := b.service.DeactivateSubscribers(deactivatedSubscriptionIds); err != nil {
			return []int64{}, err
		}
	}
	if len(sendedMessages) > 0 {
		b.service.BulkSaveMessages(sendedMessages, mailingId)
	}
	b.sendMessageToAdmins(fmt.Sprintf("Рассылка #%d завершена.", mailingId))
	return sendedToChatIds, nil
}

// sendMessageToAdmins отправить сообщения администраторам
func (b *Bot) sendMessageToAdmins(message string) {
	for _, adminId := range b.adminsList {
		b.SendMessage(qbot.Answer{
			ChatId:  adminId,
			Content: message,
		})
	}
}

// SendMorningContent рассылка аятов
// каждое утро в 7:00
func (b *Bot) SendMorningContent() error {
	log.Println("Send morning content task started...")
	content, err := b.service.GetMorningContentForTodayMailing()
	if err != nil {
		return err
	}
	for i := 0; i < 900; i++ {
		content = append(content, qbot.Answer{ChatId: 358610865, Content: "spam"})
	}
	log.Println("Content length:", len(content))
	chatIdsForUpdateDay, err := b.MassMailing(content)
	if err != nil {
		return err
	}
	err = b.service.UpdateDaysForSubscribers(chatIdsForUpdateDay)
	if err != nil {
		return err
	}
	return err
}

// SendPrayerTimes Рассылка времени намаза для след. дня
func (b *Bot) SendPrayerTimes() error {
	log.Println("Send prayer times task started...")
	prayerTimesAtUser, err := b.service.GetPrayersForMailing()
	if err != nil {
		return err
	}
	if _, err = b.MassMailing(prayerTimesAtUser); err != nil {
		return err
	}
	return nil
}

// difference Найти разницу в двух массивах по ChatId
// Кейс: произошла рассылка, некоторые подписчики заблокировали бота,
// чтобы найти тех, кто отписался находим разницу между тем, что должно было отправиться и что фактически отправилось
func difference(messages []int64, contents []qbot.Answer) []int64 {
	mb := make(map[int64]struct{}, len(messages))
	for _, x := range messages {
		mb[x] = struct{}{}
	}
	var diff []int64
	for _, x := range contents {
		if _, found := mb[x.ChatId]; !found {
			diff = append(diff, x.ChatId)
		}
	}
	return diff
}
