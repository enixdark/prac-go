package simple

import (
	"bufio"
	"crypto/sha1"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
	"text/template"
	"time"
)

var WelcomeTitle = "tou've successfully registered"
var WelcomeEmail = "Welcome to our CMS, {{Email}}! We're glad you cound join us."

const (

)