package main

import (
	"feklistova/models"
	"log"
	"net/http"
	"os"

	"context"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	// Retrieve form variables
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := repo.GetUserByEmail(ctx, email)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if password != user.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, "secret")
	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")
	if password != confirmPassword {
		http.Error(w, "password doesn't match confirm password", http.StatusBadRequest)
		return
	}

	log.Printf("Registering user: %s, %s, %s", name, email, password)

	var userID int
	userID, err = repo.RegisterUser(ctx, models.User{
		Username:  name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	})
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	log.Printf("User %s registered successfully with ID: %d", name, userID)

	// создаем сессию
	session, _ := store.Get(r, "secret")
	session.Values["authenticated"] = true
	session.Values["user_id"] = userID
	session.Save(r, w)

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func IsAuthorized(r *http.Request) bool {
	session, _ := store.Get(r, "secret")
	authenticated := session.Values["authenticated"]
	return authenticated != nil && authenticated.(bool)
}

func GetUserID(r *http.Request) int {
	session, _ := store.Get(r, "secret")
	userId := session.Values["user_id"]
	intValue, ok := userId.(int)
	if !ok {
		log.Println("Error while getting user id - not int")
		return -1
	}
	return intValue
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if IsAuthorized(r) {
		htmlData, err := os.ReadFile("web/profile.html")
		if err != nil {
			http.Error(w, "Error reading HTML file", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err = w.Write(htmlData)
		if err != nil {
			log.Println("Error serving HTML file:", err)
		}
	} else {
		http.Redirect(w, r, "/users/enter", http.StatusSeeOther)
	}
}
