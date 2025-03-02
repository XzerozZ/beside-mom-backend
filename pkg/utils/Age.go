package utils

import (
	"fmt"
	"time"
)

func CalculateAgeInDays(birthDate time.Time) (int, error) {
	now := time.Now()
	if birthDate.After(now) {
		return 0, fmt.Errorf("birthdate cannot be in the future")
	}

	duration := now.Sub(birthDate)
	days := int(duration.Hours() / 24)
	return days, nil
}

func CalculateAgeInMonths(birthDate time.Time) (int, error) {
	now := time.Now()
	if birthDate.After(now) {
		return 0, fmt.Errorf("birthdate cannot be in the future")
	}

	months := (now.Year() - birthDate.Year()) * 12
	months += int(now.Month() - birthDate.Month())
	if now.Day() < birthDate.Day() {
		months--
	}

	return months, nil
}

func CalculateAge(birthDate time.Time) (int, error) {
	now := time.Now()
	if birthDate.After(now) {
		return 0, fmt.Errorf("birthdate cannot be in the future")
	}

	years := now.Year() - birthDate.Year()
	if now.Month() < birthDate.Month() ||
		(now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		years--
	}

	return years, nil
}
