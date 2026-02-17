package models

type AladhanResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		Timings Timings `json:"timings"`
		Date    Date    `json:"date"`
	} `json:"data"`
}

type Timings struct {
	Fajr    string `json:"Fajr"`
	Dhuhr   string `json:"Dhuhr"`
	Asr     string `json:"Asr"`
	Maghrib string `json:"Maghrib"`
	Isha    string `json:"Isha"`
}

type Date struct {
	Hijri Hijri `json:"hijri"`
}

type Hijri struct {
	Day   string `json:"day"`
	Month struct {
		En string `json:"en"`
	} `json:"month"`
	Year string `json:"year"`
}

type PrayerTimes struct {
	Fajr    string
	Dhuhr   string
	Asr     string
	Maghrib string
	Isha    string
}
