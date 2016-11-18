package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	templates = template.Must(template.ParseGlob("templates/*"))
	port = ":8000"
)

func startHandler(w http.ResponseWriter, r *http.Request){
	err := templates.ExecuteTemplate(w, "flash.html", nil)
	if err != nil {
		log.Fatal("Template flash missing")
	}
}

func middleHandler(w http.ResponseWriter, r * http.Request){
	cookieValue := r.PostFormValue("message")
	cookie := http.Cookie{Name: "message",
		Value: "message:" + cookieValue,
		Expires: time.Now().Add(60 * time.Second),
		HttpOnly: true}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/finish", 301)
}

func finishHandler(w http.ResponseWriter, r* http.Request){
	cookieVal, _ := r.Cookie("message")
	if cookieVal != nil {
		fmt.Fprintln(w, "We found: " + string(cookieVal.Value) + ", but try to refresh")
		cookie := http.Cookie{Name: "message", Value: "", Expires: time.Now(), HttpOnly: true}
		http.SetCookie(w, &cookie)
	} else {
		fmt.Fprintln(w, "That cookie was gone in a flash")
	}
}

func Render(w http.ResponseWriter, r* http.Request){
	t, _ := template.ParseFiles("templates/handler.html")
	t.Execute(w, nil)
}

func main(){
	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/middle", middleHandler)
	http.HandleFunc("/finish", finishHandler)
	http.HandleFunc("/", Render)
	log.Fatal(http.ListenAndServe(port, nil))
}



