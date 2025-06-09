package utils

import (
	"fmt"
	"time"
)

func CalculateAgeDetailed(birthDate time.Time) (int, int, int, error) {
	now := time.Now()
	if birthDate.After(now) {
		return 0, 0, 0, fmt.Errorf("birthdate cannot be in the future")
	}

	years := now.Year() - birthDate.Year()
	months := int(now.Month()) - int(birthDate.Month())
	days := now.Day() - birthDate.Day()

	if days < 0 {
		months--
		lastMonth := now.AddDate(0, -1, 0)
		daysInLastMonth := time.Date(lastMonth.Year(), lastMonth.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()
		days += daysInLastMonth
	}

	if months < 0 {
		years--
		months += 12
	}

	return years, months, days, nil
}

func CalculateAgeAdjusted(birthDate time.Time, beforeBirth int) (int, int, int, error) {
	now := time.Now()
	adjustmentDays := (40 * 7) - (beforeBirth * 7)
	adjustedDate := birthDate.AddDate(0, 0, adjustmentDays)
	if adjustedDate.After(now) {
		return 0, 0, 0, fmt.Errorf("adjusted date cannot be in the future")
	}

	years := now.Year() - adjustedDate.Year()
	months := int(now.Month()) - int(adjustedDate.Month())
	days := now.Day() - adjustedDate.Day()

	if days < 0 {
		months--
		lastMonth := now.AddDate(0, -1, 0)
		daysInLastMonth := time.Date(lastMonth.Year(), lastMonth.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()
		days += daysInLastMonth
	}

	if months < 0 {
		years--
		months += 12
	}

	return years, months, days, nil
}

func CompareAgeKid(birthDate time.Time, date time.Time) (int, error) {
	months := (date.Year() - birthDate.Year()) * 12
	months += int(date.Month() - birthDate.Month())
	if date.Day() < birthDate.Day() {
		months--
	}

	return months, nil
}
