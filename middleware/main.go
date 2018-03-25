package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"strconv"
	"time"
	"github.com/justinas/alice"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"os"
)

type city struct {
	Name string
	Area uint64
}

func middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Excecuting middleware before request phase!")
		handler.ServeHTTP(w, r)
		fmt.Println("Executing middleware after response phase!")
	})
}

func filterContentType(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("415 - Unsupported Media Type. Please send JSON"))
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func setServerTimeCookie(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		cookie := http.Cookie{Name: "Server-Time(UTC)", Value: strconv.FormatInt(time.Now().Unix(), 10)}
		http.SetCookie(w, &cookie)
		log.Println("Currently in the set server time middleware")
	})
}

func mainLogic(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing request!")
	if r.Method == "POST" {
		var tempCity city
		decorder := json.NewDecoder(r.Body)
		err := decorder.Decode(&tempCity)
		if err != nil {
			panic(err)
		}

		defer r.Body.Close()

		fmt.Printf("Got %s city with area of %d sq miles!\n",
			tempCity.Name, tempCity.Area)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("201 - Created"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte("405 - Method Not Allowed"))
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainLogic)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	mainLogicHandler := http.HandlerFunc(mainLogic)
	chain := alice.New(filterContentType, setServerTimeCookie).Then(mainLogicHandler)
	http.Handle("/", middleware(mainLogicHandler))
	http.Handle("/city", chain)
	
	http.ListenAndServe(":8000", loggedRouter)
}