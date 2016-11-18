package main

import (
  "database/sql"
  "fmt"
  _ "github.com/go-sql-driver/mysql"
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "html/template"
  "k8s.io/kubernetes/pkg/util/json"
  //"crypto/tls"
  "strconv"
  "time"
  "regexp"
  "io"
  "github.com/gorilla/sessions"
  "crypto/rand"
  "crypto/sha1"
  "encoding/base64"
  "github.com/streadway/amqp"
  "text/template"
  "bytes"
)

const (
  DBHost = "localhost:3306"
  DBPort = ":3306"
  DBUser = "root"
  DBPass = "root"
  DBDatabase = "go"
  Port = ":8000"
  MQHost = "127.0.0.1"
  MQPort = ":5672"
  MQUser = "guest"
  MQPass = "guest"
)

var database *sql.DB
var sessionStore = sessions.NewCookieStore([]byte("our-social-network-apllication"))
var UserSession Session
var WelcomeTitle = "You've successfully registered!"
var WelcomeEmail = "Welcome to our CMS, {{Email}}! We're glad you could join us."


type Page struct{
  Id int
  Title string
  RawContent string
  Content template.HTML
  Date string
  Guid string
  Comments []Comment
  Session Session
}

type JSONResponse struct{
  Fields map[string]string
}

type Comment struct{
  Id int
  Name string
  Email string
  CommentText string
}

type Cookie struct{
  Name string
  Value string
  Path string
  Domain string
  Expires time.Time
  RawExpires string
  MaxAge int
  Secure bool
  HttpOnly bool
  Raw string
  Unparsed []string
}

type User struct{
  Id int
  Name string
  Password string
}

type Session struct{
  Id string
  Authenticated bool
  Unauthenticated bool
  User User
}

type RegistrationData struct{
  Email string `json:email`
  Message string `json:message`
}

var (
  templates = template.Must(template.ParseGlob("templates/*"))
  port = ":8000"
)

func MQConnect() (*amqp.Connection, *amqp.Channel, error) {
  url := "amqp://" + MQUser + ":" + MQPass + "@" + MQUser + MQPort
  conn, err := amqp.Dial(url)
  if err != nil {
    return nil, nil, err
  }
  channel, err := conn.Channel()

  if err != nil {
    return nil, nil, err
  }

  if _, err := channel.QueueDeclare("", false, true, false, false, nil); err != {
    return nil, nil, err
  }
  return conn, channel, nil
}

func MQPublish(message []byte){
  err = channel.Publish(
    "email",
    "",
    false,
    false,
    amqp.Publishing{
      ContentType: "text/plain",
      Body: []byte(message),
    }
  )
}

//func startHandler(w http.ResponseWriter, r *http.Request){
//  err := templates.ExecuteTemplate(w, "flash.html", nil)
//  if err != nil {
//    log.Fatal("Template flash missing")
//  }
//}
//
//func middleHandler(w http.ResponseWriter, r * http.Request){
//  cookieValue := r.PostFormValue("message")
//  cookie := http.Cookie{Name: "message",
//    Value: "message:" + cookieValue,
//    Expires: time.Now().Add(60 * time.Second),
//    HttpOnly: true}
//  http.SetCookie(w, &cookie)
//  http.Redirect(w, r, "/finish", 301)
//}


func getSessionUID(sid string) int{
  user := User{}
  err := database.QueryRow("Select user_id from sessions where session_id=?", sid).Scan(&user.Id)
  if err != nil {
    fmt.Println(err.Error)
    return 0
  }
  return user.Id
}

func updateSession(sid string, uid int){
  const timeFmt = "2006-01-02T15:04:05.999999999"
  timestamp := time.Now().Format(timeFmt)
  _, err := database.Exec("Insert into sessions Set session_id=?, user_id=?,session_update=? " +
          "on duplicate key update user_id=?, session_update=?", sid, uid, timestamp, uid, timestamp)
  if err != nil {
    fmt.Println(err.Error)
  }
}

func generateSessionId() string{
  sid := make([]byte, 24)
  _, err := io.ReadFull(rand.Reader, sid)
  if err != nil {
    log.Fatal("Could not generate session id")
  }
  return base64.URLEncoding.EncodeToString(sid)
}

func validateSession(w http.ResponseWriter, r *http.Request){
  session, _ := sessionStore.Get(r, "app-session")
  if sid, valid := session.Values["sid"]; valid {
    currentUID := getSessionUID(sid.(string))
    updateSession(sid.(string), currentUID)
    UserSession.Id = string(currentUID)
  } else {
    newSID := generateSessionId()
    session.Values["sid"] = newSID
    session.Save(r, w)
    UserSession.Id = newSID
    updateSession(newSID, 0)
  }
  fmt.Println(session.ID)
}

func ServerPage(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  pageId := vars["guid"]
  thisPage := Page{}
  fmt.Println(pageId)

  err := database.QueryRow("Select id, page_title, page_content, " +
          "page_date, page_guid from pages where page_guid=?", pageId).
          Scan(&thisPage.Id, &thisPage.Title, &thisPage.RawContent, &thisPage.Date, &thisPage.Guid)

  if err != nil {
    http.Error(w, http.StatusText(404), http.StatusNotFound)
    log.Println("Couldn't get page: " + pageId)
    log.Println(err.Error)
    return
  }

  comments , err := database.Query("select id, comment_name as Name, comment_email," +
          "comment_text from comments where page_id=?", thisPage.Id )
  if err != nil {
    log.Println(err.Error)
  }

  for comments.Next() {
    var comment Comment

    comments.Scan(&comment.Id, &comment.Name, &comment.Email, &comment.CommentText)
    thisPage.Comments = append(thisPage.Comments, comment)
  }
  thisPage.Content = template.HTML(thisPage.RawContent)
  //thisPage.GUID = pageId
  t, _ := template.ParseFiles("templates/blog.html")

  // html := `<html><head><title>` + thisPage.Title +
  // `</title></head><body><h1>` + thisPage.Title + `</h1><div>` +
  // thisPage.Content + `</div></body></html>`
  t.Execute(w, thisPage)
  // fmt.Fprintln(w, html)
}

func RedirIndex(w http.ResponseWriter, r *http.Request){
  http.Redirect(w, r, "/home", 301)
}

func ServerIndex(w http.ResponseWriter, r *http.Request){
  var Pages = []Page{}
  pages, err := database.Query("Select page_title, page_content, page_date, page_guid from pages order by ? DESC",
    "page_date")
  if err != nil {
    fmt.Fprint(w, err.Error)
  }

  defer pages.Close()
  for pages.Next() {
    thisPage := Page{}
    pages.Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date, &thisPage.Guid)
    thisPage.Content = template.HTML(thisPage.RawContent)
    Pages = append(Pages, thisPage)
  }

  t, _ := template.ParseFiles("templates/index.html")
  t.Execute(w, Pages)
}

func APIPage(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  pageGUID := vars["guid"]
  thisPage := Page{}

  fmt.Println(pageGUID)

  err := database.QueryRow(" select page_title, page_content, page_date from pages where page_guid=?", pageGUID).
          Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)

  thisPage.Content = template.HTML(thisPage.RawContent)

  if err != nil {
    http.Error(w, http.StatusText(404), http.StatusNotFound)
    log.Println(err)
    return
  }

  APIOutput, err := json.Marshal(thisPage)
  fmt.Println(APIOutput)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintln(w, thisPage)

}


func APICommentPost(w http.ResponseWriter, r *http.Request){
  var commentAdded bool
  err := r.ParseForm()
  if err != nil {
    log.Println(err.Error)
  }

  name := r.FormValue("name")
  email := r.FormValue("email")
  comments := r.FormValue("comments")

  res, err := database.Exec("INSERT INTO comments SET comment_name =?, " +
          "comment_email=?, comment_text=?", name, email, comments)
  if err != nil {
    log.Println(err.Error)
  }

  id, err := res.LastInsertId()
  if err != nil {
    commentAdded = false
  } else{
    commentAdded = true
  }

  commentAddedBool := strconv.FormatBool(commentAdded)

  var resp JSONResponse
  resp.Fields = make(map[string]string)
  resp.Fields["id"] = string(id)
  resp.Fields["added"] = commentAddedBool

  jsonResp, _ := json.Marshal(resp)
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintln(w, jsonResp)

}

func APICommentPut(w http.ResponseWriter, r *http.Request){
  err := r.ParseForm()
  if err != nil {
    log.Println(err.Error)
  }
  vars := mux.Vars(r)
  id := vars["id"]
  fmt.Println(id)

  name := r.FormValue("name")
  email := r.FormValue("email")
  comments := r.FormValue("comments")

  res, err := database.Exec("Update comments set comment_name=?, comment_email=?, " +
          "comment_text=? where comment_id=?", name, email, comments, id)
  fmt.Println(res)
  if err != nil {
    log.Println(err.Error)
  }

  var resp JSONResponse

  jsonResp, _ := json.Marshal(resp)
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintln(w, jsonResp)
}

func weakPasswordHash(password string) []byte{
  hash := sha1.New()
  io.WriteString(hash, password)
  return hash.Sum(nil)

}

func RegisterGet(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("templates/register.html")
  t.Execute(w, nil)
}

func RegisterPost(w http.ResponseWriter, r *http.Request){
  err := r.ParseForm()
  if err != nil {
    log.Fatal(err.Error)
  }

  name := r.FormValue("user_name")
  email := r.FormValue("user_email")
  pass := r.FormValue("user_password")
  //pass_confirmation := r.FormValue("user_password_confirmation")
  pageGUID := r.FormValue("referrer")
  gure := regexp.MustCompile("[^A-Za-z0-9]+")
  guid := gure.ReplaceAllString(name, "")
  password := weakPasswordHash(pass)

  res, err := database.Exec("Insert into users set user_name=?, user_guid=?" +
          "user_email=?, user_password=?", name, guid, email, password)

  fmt.Println(res)

  if err != nil {
    fmt.Fprintln(w, err.Error)
  } else {
    Email := RegistrationData{Email:: email, Message: ""}
    messsage, err := template.New("email").Parse(WelcomeEmail)
    var mbuf bytes.Buffer
    messsage.Execute(&mbuf, Email)
    MQPublish(json.Marshal(mbuf.String()))
    http.Redirect(w, r, "/page/" + pageGUID, 301)
  }
}

func LoginGet(w http.ResponseWriter, r* http.Request){
  t, _ := template.ParseFiles("templates/login.html")
  t.Execute(w, nil)
}

func LoginPost(w http.ResponseWriter, r *http.Request){

  //thisUser := User{}
  //err := r.ParseForm()
  //if err != nil {
  //  log.Fatal(err.Error)
  //}
  //
  //name := r.FormValue("user_name")
  //pass := r.FormValue("user_password")
  //
  //err = database.QueryRow("select user_name ,user_password " +
  //        "from users where user_name = ?", name).
  //        Scan(&thisUser.Name, &thisUser.Password)
  //
  //if err != nil {
  //  fmt.Fprintln(w, err.Error)
  //}
  //
  //password := weakPasswordHash(pass)
  //if password == thisUser.Password {
  //  http.Redirect(w, r , "/", 301)
  //} else{
  //  fmt.Fprintln(w, err.Error)
  //}

  validateSession(w,r)
  user := User{}
  name := r.FormValue("user_name")
  pass := r.FormValue("user_password")
  password := weakPasswordHash(pass)
  err := database.QueryRow("Select user_id, user_name, " +
          "from users where user_name=? and user_password=?", name, password).Scan(&user.Id, &user.Name)

  if err != nil {
    fmt.Fprintln(w, err.Error)
    user.Id = 0
    user.Name = ""
  } else {
    updateSession(UserSession.Id, user.Id)
    fmt.Fprintln(w, user.Name)
  }
}



func main(){
  dbConn := fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPass, DBHost, DBDatabase)
  fmt.Println(dbConn)
  db, err := sql.Open("mysql", dbConn)
  if err != nil {
    log.Println("Couldn't connect!")
    log.Println(err.Error)
  }

  database = db
  routes := mux.NewRouter()
  //routes.HandleFunc("/api/pages", APIPage).Methods("GET").Schemes("https")
  //routes.HandleFunc("/api/pages/{guid:[0-9a-zA\\-]+}", APIPage).Methods("GET").Schemes("https")
  routes.HandleFunc("/api/comments", APICommentPost).Methods("POST")
  routes.HandleFunc("/api/comments/{id:[\\w\\d\\-]+}", APICommentPut).Methods("PUT")
  routes.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", ServerPage)
  routes.HandleFunc("/register", RegisterGet).Methods("GET")
  routes.HandleFunc("/register", RegisterPost).Methods("POST")
  routes.HandleFunc("/login", LoginPost).Methods("POST")
  routes.HandleFunc("/login", LoginGet).Methods("GET")
  routes.HandleFunc("/", RedirIndex)
  routes.HandleFunc("/home", ServerIndex)
  http.Handle("/", routes)
  http.ListenAndServe(Port, nil)
  //certificates, err := tls.LoadX509KeyPair("ssl/cert.pem", "ssl/key.pem")
  //tlsConf := tls.Config{Certificates:
  //[]tls.Certificate{certificates}}
  //tls.Listen("tcp", Port, &tlsConf)

}
