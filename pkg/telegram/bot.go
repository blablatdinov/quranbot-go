package telegram

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"qbot"
	"qbot/pkg/service"

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
		log.Printf("Bot started on long polling... Debug: %v", b.bot.Debug)
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

// MassMailing функция для массовой рассылки
func (b *Bot) MassMailing(content []qbot.Answer) ([]int64, error) {
	if len(content) == 0 {
		return []int64{}, errors.New("content len == 0")
	}
	mailingId, err := b.service.CreateMailing()
	if err != nil {
		return []int64{}, err
	}
	var sendedMessages []tgbotapi.Message
	var deactivatedSubscriptionIds []int64
	var sendedToChatIds []int64
	log.Printf("Start mailing (%d)...", mailingId)
	for _, _elem := range content {
		sliceIndex := minNumber(len(_elem.Content), 50)
		log.Printf("Send message to %d: (%s)...", _elem.ChatId, _elem.Content[:sliceIndex])
		message, err := b.SendMessage(_elem)
		if err != nil {
			if err.Error() == "FIXME" {
				deactivatedSubscriptionIds = append(deactivatedSubscriptionIds, _elem.ChatId)
			}
		}
		sendedMessages = append(sendedMessages, message)
		sendedToChatIds = append(sendedToChatIds, message.Chat.ID)
	}
	log.Printf("Mailing (%d) sended", mailingId)
	if len(sendedMessages) > 0 {
		b.service.BulkSaveMessages(sendedMessages, mailingId)
	}
	if len(deactivatedSubscriptionIds) > 0 {
		if err := b.service.DeactivateSubscribers(deactivatedSubscriptionIds); err != nil {
			return []int64{}, err
		}
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
