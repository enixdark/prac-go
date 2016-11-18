package main

import (
  "net/http"
  "github.com/gorilla/mux"
  "fmt"
  "os"
)

const ( Port = ":8000" )

func pageHandler(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  pageID := vars["id"]
  fmt.Println(pageID)
  fileName := "files/" + pageID + ".html"
  _, err := os.Stat(fileName)
  if err != nil {
    fileName = "files/404.html"
  }
  http.ServeFile(w, r, fileName)
}

func main(){
  
  rtr := mux.NewRouter()
  rtr.HandleFunc("/pages/{id:[0-9]+}", pageHandler)
  //rtr.HandleFunc("/homepage", pageHandler)
  //rtr.HandleFunc("/contact", pageHandler)
  http.Handle("/", rtr)
  http.ListenAndServe(Port,nil) 
}
