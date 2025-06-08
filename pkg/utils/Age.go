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
	if now.Month() < birthDate.Month() ||
		(now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		years--
	}

	currentYearBirthday := time.Date(now.Year(), birthDate.Month(), birthDate.Day(), 0, 0, 0, 0, time.Local)
	if now.Before(currentYearBirthday) {
		currentYearBirthday = currentYearBirthday.AddDate(-1, 0, 0)
	}

	daysAfterLastBirthday := int(now.Sub(currentYearBirthday).Hours() / 24)
	weeks := daysAfterLastBirthday / 7
	days := daysAfterLastBirthday % 7

	return years, weeks, days, nil
}

func CompareAgeKid(birthDate time.Time, date time.Time) (int, error) {
	months := (date.Year() - birthDate.Year()) * 12
	months += int(date.Month() - birthDate.Month())
	if date.Day() < birthDate.Day() {
		months--
	}

	return months, nil
}
