package main

import (
	"math/rand"
	"fmt"
	"net/http"
	"os/exec"
	"bytes"
	"log"
	"github.com/gorilla/mux"
	"time"
)

func giveRandom(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Your random number is: %f", rand.Float64())
}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category is: %v\n", vars["category"])
    fmt.Fprintf(w, "ID is: %v\n", vars["id"])
}

func getCommandOutput(command string, arguments ...string) string {
	cmd := exec.Command(command, arguments...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
	}

	err = cmd.Wait()
    if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
	}

	return out.String()
}

func goVersion(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, getCommandOutput("go", "version"))
}

func getFileContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, getCommandOutput("/bin/cat", vars["name"]))
}

func randomFloat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, rand.Float64())
}

func randomInt(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, rand.Int())
}

func main() {

	
	router := mux.NewRouter()

	// router.ServeFiles("/static/*filepath", http.Dir("../static"))
	
	router.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
	router.HandleFunc("/api/v1/go-version", goVersion)
    router.HandleFunc("/api/v1/show-file/{name}", getFileContent)
	router.HandleFunc("/randomFloat", randomFloat)
	router.HandleFunc("/randomInt", randomInt)

	srv := &http.Server{
		Handler: router,
		Addr: "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

	

}

