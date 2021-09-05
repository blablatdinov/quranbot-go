package telegram

import "log"

func (b *Bot) StartJobs() error {
	jobs := map[string]interface{}{
		"0 7 * * *": b.SendMorningContent,
	}
	for cronTime, job := range jobs {
		_, err := b.goCron.Cron(cronTime).Do(job.(func() error))
		if err != nil {
			return err
		}
	}
	b.goCron.StartAsync()
	log.Printf("Cron started with jobs: %s\n", jobs)
	return nil
}
