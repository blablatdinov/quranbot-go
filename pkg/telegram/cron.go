package telegram

import (
	"fmt"
	"log"
)

func (b *Bot) StartJobs() error {
	jobs := map[string]interface{}{
		"0 7 * * *": b.SendMorningContent,
		// "0 20 * * *": b.SendPrayerTimes,
	}
	for cronTime, job := range jobs {
		_, err := b.goCron.Cron(cronTime).Do(job.(func() error))
		if err != nil {
			return err
		}
	}
	b.goCron.StartAsync()
	logStartedJobs(jobs)
	return nil
}

func logStartedJobs(jobs map[string]interface{}) {
	message := "Cron started with jobs:\n"
	for time, job := range jobs {
		message = message + fmt.Sprintf("\t%s %s\n", time, job)
	}
	log.Printf(message, jobs)
}
