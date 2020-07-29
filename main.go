package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/onurcevik/postgresboilerplate/database"
	"github.com/onurcevik/postgresboilerplate/handlers"
)

var (
	host     = os.Getenv("HOST")
	port, _  = strconv.Atoi(os.Getenv("POSTGRESPORT"))
	user     = os.Getenv("USER")
	password = os.Getenv("PASSWORD")
	dbname   = os.Getenv("DBNAME")
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	database.Conn, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer database.Conn.Close()
	err = database.Conn.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Database")

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/dashboard", handlers.DashboardHandler)
	http.ListenAndServe(":8080", nil)
}
