package handlers

import (
	"html/template"
	"net/http"
)

func IndexHandler(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("html/index.html")
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(writer, nil)
	if err != nil {
		return
	}
}
