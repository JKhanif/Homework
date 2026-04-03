package service

import (
	"context"
	"fmt"
	"log"
	"prayertimes/clients/aladhan"
	"time"
)

type Service struct {
	aladhanClient *aladhan.Client
}

func New(aladhanClient *aladhan.Client) *Service {
	return &Service{
		aladhanClient: aladhanClient,
	}
}

var prayerNames = [5]string{"Фаджр", "Зухр", "Аср", "Магриб", "Иша"}

func (s *Service) ResponsePrayerTime(ctx context.Context, city string) (string, error) {
	res, err := s.aladhanClient.GetTodayPrayerTimesByCity(ctx, city)
	if err != nil {
		log.Println("error s.aladhanClient.GetTodayPrayerTimesByCity: ", err)
		return "", err
	}

	timeStr := []string{res.Data.Timings.Fajr, res.Data.Timings.Dhuhr, res.Data.Timings.Asr, res.Data.Timings.Maghrib, res.Data.Timings.Isha}

	nextPrayerFound := false
	resMessage := ""
	now := time.Now()
	//currentIndex := -1

	for i, v := range timeStr {
		parsed, err := time.Parse("15:04", v)
		if err != nil {
			log.Println("error time.Parse: ", err)
			return "", err
		}

		prayerTime := time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), 0, 0, now.Location())

		// if !now.Before(prayerTime) {
		// 	currentIndex = i
		// }

		if nextPrayerFound == false && now.Before(prayerTime) {
			nextPrayerFound = true
			nextTP := prayerTime.Sub(now)
			h := int(nextTP.Hours())
			m := int(nextTP.Minutes()) % 60
			timeLeft := fmt.Sprintf("%02d:%02d", h, m)
			resMessage += fmt.Sprintf("<b>%s: %s</b>", prayerNames[i], v)
			resMessage += fmt.Sprintf("(Через %s)\n", timeLeft)
			continue
		}

		resMessage += fmt.Sprintf("%s: %s\n", prayerNames[i], v)
	}

	format := time.Now().Format("02.01.2006")

	hijriTime := fmt.Sprintf("\n%s %s %s", res.Data.Date.Hijri.Day, res.Data.Date.Hijri.Month.En, res.Data.Date.Hijri.Year)
	resMessage = fmt.Sprintf("🕌Город: %s\nРасписание на %s\n\n%s%s", city, format, resMessage, hijriTime)

	return resMessage, nil
}
