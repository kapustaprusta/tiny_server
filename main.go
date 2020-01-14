package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	db *sql.DB
)

type API struct {
	Message string `json:"message"`
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
	First string `json:"first"`
	Last  string `json:"last"`
}

type Users struct {
	Users []User `json:"users"`
}

type CreateResponse struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"code"`
}

type ErrorMessage struct {
	Code    int
	Status  int
	Message string
}

func init() {
	db, _ = sql.Open("mysql", "vladislav:password@/social_network")
}

func dbErrorParse(err string) (string, int64) {
	errorParts := strings.Split(err, ":")
	errorMessage := errorParts[1]
	errorCode, _ := strconv.ParseInt(strings.Split(errorParts[0], "Error ")[1], 10, 32)

	return errorMessage, errorCode

}

func translateError(err int64) ErrorMessage {
	var errMsg ErrorMessage
	switch err {
	case 1062:
		errMsg.Message = "Duplicate entry"
		errMsg.Code = 10
		errMsg.Status = 409
	}

	return errMsg
}

func CreateUsers(w http.ResponseWriter, r *http.Request) {
	newUser := User{}
	newUser.Name = r.FormValue("user")
	newUser.Email = r.FormValue("email")
	newUser.First = r.FormValue("first")
	newUser.Last = r.FormValue("last")

	_, err := json.Marshal(newUser)
	if err != nil {
		log.Println("Something went wrong! ", err)
		return
	}

	response := CreateResponse{}
	sql := "INSERT INTO users set user_name='" + newUser.Name + "', user_first='" + newUser.First + "', user_last='" + newUser.Last + "', user_email='" + newUser.Email + "'"
	_, err = db.Exec(sql)
	if err != nil {
		errorMessage, errorCode := dbErrorParse(err.Error())
		translatingError := translateError(errorCode)

		response.Error = translatingError.Message
		response.ErrorCode = translatingError.Code

		log.Println(errorMessage, translatingError.Status)
		http.Error(w, "Conflict", translatingError.Status)
	}

	output, _ := json.Marshal(response)
	fmt.Fprintf(w, string(output))
}

func RetrieveUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("starting retrieval")
	start, limit := 0, 10
	next := start + limit

	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Link", "<http://localhost:8080/api/users?start="+string(next)+"; rel=\"next\"")

	rows, _ := db.Query("SELECT * FROM users LIMIT 10")
	response := Users{}

	for rows.Next() {
		user := User{}
		rows.Scan(&user.ID, &user.Name, &user.First, &user.Last, &user.Email)
		response.Users = append(response.Users, user)
	}

	output, _ := json.Marshal(response)
	fmt.Fprintf(w, string(output))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/users", CreateUsers).Methods("POST")
	router.HandleFunc("/api/users", RetrieveUsers).Methods("GET")

	log.Println("Server was started...")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
