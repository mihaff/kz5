package main

import (
	"io"
	"net/http"
	"os"
)

func handlerTmpl(w http.ResponseWriter, r *http.Request, template_file string) {
	file, err := os.Open(template_file)
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/html")

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
		return
	}
}

func HomeHandlerTmpl(w http.ResponseWriter, r *http.Request) {
	handlerTmpl(w, r, "web/index.html")
}

func LoginHandlerTmpl(w http.ResponseWriter, r *http.Request) {
	handlerTmpl(w, r, "web/enter.html")
}

func RegisterHandlerTmpl(w http.ResponseWriter, r *http.Request) {
	handlerTmpl(w, r, "web/registration.html")
}

func ProfileHandlerTmpl(w http.ResponseWriter, r *http.Request) {
	handlerTmpl(w, r, "web/profile.html")
}

func ProgressClassHandlerTmpl(w http.ResponseWriter, r *http.Request) {
	if IsAuthorized(r) {
		handlerTmpl(w, r, "web/model_form_class.html")
	} else {
		http.Redirect(w, r, "/users/enter", http.StatusSeeOther)
	}
}

func ProgressRegHandlerTmpl(w http.ResponseWriter, r *http.Request) {
	if IsAuthorized(r) {
		handlerTmpl(w, r, "web/model_form_reg.html")
	} else {
		http.Redirect(w, r, "/users/enter", http.StatusSeeOther)
	}
}

func ResultClassHandlerTmpl(w http.ResponseWriter, r *http.Request) {
	handlerTmpl(w, r, "web/save_model_class.html")
}

func ResultRegHandlerTmpl(w http.ResponseWriter, r *http.Request) {
	handlerTmpl(w, r, "web/save_model_reg.html")
}
