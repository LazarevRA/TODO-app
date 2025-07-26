package api

import (
	"net/http"
)

const Layout = "20060102"

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", TaskHandler)
}
