package api

import (
	"net/http"
	"time"
)

const Layout = "20060102"

var NowTest = time.Date(2024, time.January, 26, 0, 0, 0, 0, time.UTC)

func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", TaskHandler)
}
