package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/onurcevik/postgresboilerplate/database"
)

func getUserFromServer(req *http.Request) string {

	c, err := req.Cookie("session")
	if err != nil {
		panic(err)
	}
	sqlstmnt := `SELECT sessions.username FROM sessions INNER JOIN users ON sessions.username=users.username WHERE sessions.value=$1`
	row := database.Conn.QueryRow(sqlstmnt, c.Value)

	var username string
	switch err := row.Scan(&username); err {
	case sql.ErrNoRows:
		fmt.Println("User doesnt exist in database")
	case nil:
		fmt.Println("sessions username : ", username)

	default:
		panic(err)
	}
	return username
}

func alreadyLoggedIn(req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}

	var bl bool
	sqlstmnt := `SELECT EXISTS(SELECT * FROM sessions INNER JOIN users ON sessions.username=users.username WHERE sessions.value=$1);`
	_ = database.Conn.QueryRow(sqlstmnt, c.Value).Scan(&bl)
	return bl
}
