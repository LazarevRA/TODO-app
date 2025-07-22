package nextdate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const layout = "20060102"

func afterNow(date time.Time) bool {
	return time.Now().Before(date)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	//Получаем дату
	date, err := time.Parse(layout, dstart)
	if err != nil {
		return "", errors.New("error in date parsing")
	}

	//Проверяем, что правило не пустое
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {

	case "d":
		return nextDay(date, parts)

	case "y":
		return nextYear(date, parts)

	case "w":
		return nextWeek(date, parts)
	}

}

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
		nextDate.AddDate(0, 0, interval)
		if afterNow(nextDate) {
			return nextDate.Format(layout), nil
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
		nextDate.AddDate(1, 0, 0)
		if afterNow(nextDate) {
			return nextDate.Format(layout), nil
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
		return newDate.Format(layout), nil
	}

	//Перебор дней
	for i := 0; i < 400; i++ {
		newDate = newDate.AddDate(0, 0, 1)
		if afterNow(newDate) && correctWeekDay(newDate, days) {
			return newDate.Format(layout), nil
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
