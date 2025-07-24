package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func afterNow(date time.Time) bool {

	return date.After(time.Now())
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	//Получаем дату
	date, err := time.Parse(Layout, dstart)
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

	case "m":
		return nextMonth(date, parts)
	default:
		return "", fmt.Errorf("wrong repeat date format: %s", parts[0])
	}

}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nowStr := r.FormValue("now")

	var now time.Time

	if nowStr == "" {
		now = time.Now()

		//Для тестов
		//now = time.Date(2024, time.January, 26, 0, 0, 0, 0, time.UTC)
	} else {
		var err error
		now, err = time.Parse(Layout, nowStr)
		if err != nil {
			http.Error(w, "incorrect parameter 'now'", http.StatusBadRequest)
			return
		}
	}
	nextDate, err := NextDate(now, dateStr, repeat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, nextDate)

}
