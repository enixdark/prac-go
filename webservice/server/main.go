package main

import (
	"math/rand"
	"fmt"
	"github.com/enixdark/prac-go/webservice/numerals"
	"html"
	"net/http"
	"strconv"
	"strings"
	"os/exec"
	"bytes"
	"log"
	"github.com/julienschmidt/httprouter"
	// "time"
)

type ServeMux struct {

}

func (p *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		giveRandom(w, r)
		return
	}

	http.NotFound(w, r)
	return 
}

func giveRandom(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Your random number is: %f", rand.Float64())
}



func process_number(w http.ResponseWriter, r *http.Request) {
	urlPathElements := strings.Split(r.URL.Path, "/")

	if urlPathElements[1] == "roman_number" {
		number, _ := strconv.Atoi(strings.TrimSpace(urlPathElements[2]))
		if number == 0 || number > 10 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Not Found"))
		} else {
			fmt.Fprintf(w, "%q", html.EscapeString(numerals.Numerals[number]))
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad request"))
	}
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

func goVersion(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    fmt.Fprintf(w, getCommandOutput("go", "version"))
}

func getFileContent(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, getCommandOutput("/bin/cat", params.ByName("name")))
}

func randomFloat(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintln(w, rand.Float64())
}

func randomInt(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintln(w, rand.Int())
}

func main() {
	router := httprouter.New()

    router.ServeFiles("/static/*filepath", http.Dir("../static"))
	router.GET("/api/v1/go-version", goVersion)
    router.GET("/api/v1/show-file/:name", getFileContent)

	router.GET("/randomFloat", randomFloat)

	router.GET("/randomInt", randomInt)

	log.Fatal(http.ListenAndServe(":8000", router))

}

