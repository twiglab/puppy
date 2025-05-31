package gbot

import "time"

func OpeningTime(now time.Time) (time.Time, time.Time) {
	start := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, now.Location())
	return start, end
}

func NightTime(now time.Time) (time.Time, time.Time) {
	start := time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, now.Location())
	return start, end
}

func Yestoday(now time.Time) time.Time {
	return now.Add(-1 * 24 * time.Hour)
}

func BeforWeekDay(now time.Time) time.Time {
	return now.Add(-7 * 24 * time.Hour)
}
