package controllers

import "net/http"

func AdminPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Страница администратора"))
}
