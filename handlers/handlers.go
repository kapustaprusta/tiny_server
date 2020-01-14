package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"tiny_server/defs"
	"tiny_server/utils"

	_ "github.com/go-sql-driver/mysql"
)

var (
	database *sql.DB
	api      defs.API
)

func init() {
	database, _ = sql.Open("mysql", "vladislav:1234q1234q@/social_network")

	api.Endpoints = []defs.Endpoint{
		{
			URL:         "/api/users",
			Method:      "GET",
			Description: "Return a list of users with optional parameters",
		},
		{
			URL:         "/api/user",
			Method:      "POST",
			Description: "Create a user",
		},
		{
			URL:         "/api/user/XXX",
			Method:      "PUT",
			Description: "Update a user’s information",
		},
		{
			URL:         "/api/user/XXX",
			Method:      "DELETE",
			Description: "Delete a user",
		},
		{
			URL:         "/api/connections",
			Method:      "GET",
			Description: "Return a list of connections based on users",
		},
		{
			URL:         "/api/connections",
			Method:      "POST",
			Description: "Create a connection between users",
		},
		{
			URL:         "/api/connections/XXX",
			Method:      "PUT",
			Description: "Modify a connection",
		},
		{
			URL:         "/api/connections/XXX",
			Method:      "DELETE",
			Description: "Remove a connection between users",
		},
		{
			URL:         "/api/statuses",
			Method:      "GET",
			Description: "Get a list of statuses",
		},
		{
			URL:         "/api/statuses",
			Method:      "POST",
			Description: "Create a status",
		},
		{
			URL:         "/api/statuses/XXX",
			Method:      "PUT",
			Description: "Update a status",
		},
		{
			URL:         "/api/statuses/XXX",
			Method:      "DELETE",
			Description: "Delete a status",
		},
		{
			URL:         "/api/comments",
			Method:      "GET",
			Description: "Get list of comments",
		},
		{
			URL:         "/api/comments",
			Method:      "POST",
			Description: "Create a comment",
		},
		{
			URL:         "/api/comments/XXX",
			Method:      "PUT",
			Description: "Update a comment",
		},
		{
			URL:         "/api/comments/XXX",
			Method:      "DELETE",
			Description: "Delete a comment",
		},
		{
			URL:         "/api/oauth/authorize",
			Method:      "GET",
			Description: "Returns a list of users with optional parameters",
		},
		{
			URL:         "/api/oauth/token",
			Method:      "POST",
			Description: "Creates a user",
		},
		{
			URL:         "/api/oauth/revoke",
			Method:      "PUT",
			Description: "Updates a user’s information",
		},
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := defs.User{
		Nickname:  r.FormValue("nickname"),
		Email:     r.FormValue("email"),
		FirstName: r.FormValue("firstname"),
		LastName:  r.FormValue("lastname"),
	}

	_, err := json.Marshal(newUser)
	if err != nil {
		log.Println("Something went wrong! ", err)
		return
	}

	response := defs.CreationResponse{Code: 100, Message: "User was added"}
	sqlQuery := "INSERT INTO users set user_nickname='" + newUser.Nickname + "', user_firstname='" + newUser.FirstName + "', user_lastname='" + newUser.LastName + "', user_email='" + newUser.Email + "'"
	_, err = database.Exec(sqlQuery)
	if err != nil {
		errorMessage, errorCode := utils.ParseDatabaseError(err.Error())
		translatingError := utils.TranslateError(errorCode)

		response.Message = translatingError.Comment
		response.Code = translatingError.Code

		log.Println("message: ", errorMessage, "code: ", errorCode)
		http.Error(w, "Conflict", translatingError.HTTPCode)
	}

	output, _ := json.Marshal(response)
	fmt.Fprintf(w, string(output))
}

func RetrieveUsers(w http.ResponseWriter, r *http.Request) {
	start, limit := 0, 10
	next := start + limit

	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Link", "<http://localhost:8080/api/users?start="+string(next)+"; rel=\"next\"")

	rows, _ := database.Query("SELECT * FROM users LIMIT 10")
	response := defs.Users{}

	for rows.Next() {
		user := defs.User{}
		rows.Scan(&user.ID, &user.Nickname, &user.FirstName, &user.LastName, &user.Email)
		response.Users = append(response.Users, user)
	}

	output, _ := json.Marshal(response)
	fmt.Fprintf(w, string(output))
}
