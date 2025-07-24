package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func nextDay(date time.Time, parts []string) (string, error) {
	//Проверка на правильность правила
	if len(parts) != 2 {
		return "", errors.New("wrong rule for 'd' format")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", errors.New("interval must be a number")
	}

	if interval < 1 || interval > 400 {
		return "", errors.New("selected interval must be 1-400")
	}

	nextDate := date
	//Увеличиваем время, пока не будет больше текущего
	for {
		nextDate = nextDate.AddDate(0, 0, interval)
		if afterNow(nextDate) {
			return nextDate.Format(Layout), nil
		}
	}
}

func nextYear(date time.Time, parts []string) (string, error) {
	//Проверка на правильность правила
	if len(parts) != 1 {
		return "", errors.New("wrong rule for 'y' format")
	}

	nextDate := date

	for {
		nextDate = nextDate.AddDate(1, 0, 0)
		if afterNow(nextDate) {
			return nextDate.Format(Layout), nil
		}
	}
}

func nextWeek(date time.Time, parts []string) (string, error) {
	//Проверка, что есть перечисление дней
	if len(parts) != 2 {
		return "", errors.New("wrong rule for 'w' format")
	}

	days, err := parseIntList(parts[1], 1, 7)
	if err != nil {
		return "", fmt.Errorf("invalid week days: %w", err)
	}
	if len(days) == 0 {
		return "", errors.New("empty week days list")
	}

	newDate := date

	//Проверка удовлетворяет ли стартовая дата условию "позже чем сейчас"
	if afterNow(newDate) && correctWeekDay(newDate, days) {
		return newDate.Format(Layout), nil
	}

	//Перебор дней
	for i := 0; i < 400; i++ {
		newDate = newDate.AddDate(0, 0, 1)
		if afterNow(newDate) && correctWeekDay(newDate, days) {
			return newDate.Format(Layout), nil
		}
	}

	return "", errors.New("date not found within max search period")
}

func nextMonth(date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("wrong rule for 'm' format")
	}

	days, months, err := parseMonthRule(parts)
	if err != nil {
		return "", fmt.Errorf("invalid month rule: %w", err)
	}

	newDate := date

	//Проверка удовлетворяет ли стартовая дата условию "позже чем сейчас"
	if afterNow(newDate) && correctMonthDay(newDate, days, months) {
		return newDate.Format(Layout), nil
	}

	//Перебор дней
	for i := 0; i < 10000; i++ {
		newDate = newDate.AddDate(0, 0, 1)
		if afterNow(newDate) && correctMonthDay(newDate, days, months) {
			return newDate.Format(Layout), nil
		}
	}

	return "", errors.New("date not found within max search period")
}

func parseIntList(s string, min, max int) ([]int, error) {
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))

	for _, part := range parts {
		if part == "" {
			continue
		}

		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, errors.New("not a number")
		}

		if val < min || val > max {
			return nil, fmt.Errorf("value %d out of range [%d,%d]", val, min, max)
		}

		result = append(result, val)
	}
	return result, nil
}

// Функция проверяет дату на соответсвие дням недели переданным в массиве
func correctWeekDay(date time.Time, days []int) bool {
	wday := date.Weekday()
	wdayInt := int(wday)
	if wday == time.Sunday {
		wdayInt = 7
	}
	for _, d := range days {
		if d == wdayInt {
			return true
		}
	}
	return false
}

func parseMonthRule(parts []string) (days []int, months []int, err error) {

	// Парсинг дней (-2..31)
	days, err = parseIntList(parts[1], -2, 31)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid days: %w", err)
	}

	// Парсинг месяцев
	if len(parts) > 2 {
		months, err = parseIntList(parts[2], 1, 12)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid months: %w", err)
		}
	}

	return days, months, nil
}

func correctMonthDay(date time.Time, days, months []int) bool {

	if len(months) > 0 {
		found := false
		for _, m := range months {
			if m == int(date.Month()) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Проверка дня
	lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	currentDay := date.Day()

	for _, d := range days {
		switch {
		case d > 0 && currentDay == d:
			return true
		case d == -1 && currentDay == lastDay:
			return true
		case d == -2 && currentDay == lastDay-1:
			return true
		}
	}
	return false
}
