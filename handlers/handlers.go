package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"text/template"

	"github.com/onurcevik/postgresboilerplate/database"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	tpl *template.Template
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {

	if auth := alreadyLoggedIn(r); !auth {
		http.Error(w, "Not Logged In!", http.StatusInternalServerError)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	tpl.ExecuteTemplate(w, "dashboard.html", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	var name, pwd string

	if r.Method == http.MethodPost {

		name = r.FormValue("username")
		pwd = r.FormValue("password")

	}

	if len(name) > 0 {

		var usernameExists bool
		sqlstmnt := `SELECT EXISTS(SELECT * FROM users WHERE username=$1);`
		_ = database.Conn.QueryRow(sqlstmnt, name).Scan(&usernameExists)
		fmt.Println("usernameExists : ", usernameExists)
		if usernameExists {
			http.Error(w, "Username exists", http.StatusInternalServerError)
			return
		}

		insertQuery := `INSERT INTO users (username, password)
		VALUES ($1,$2 );
		`
		registerpasswd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
		_, err := database.Conn.Exec(insertQuery, name, string(registerpasswd))
		if err != nil {
			panic(err)
		}

		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}

		http.SetCookie(w, c)

		insertStmnt := `INSERT INTO sessions (value,username) VALUES($1,$2)`
		_, err = database.Conn.Exec(insertStmnt, c.Value, name)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return

	}

	tpl.ExecuteTemplate(w, "register.html", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	//var u User
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {

		name := r.FormValue("username")
		pwd := r.FormValue("password")

		selectQuery := `SELECT password FROM users WHERE username=$1;`
		row := database.Conn.QueryRow(selectQuery, name)

		var hash string

		switch err := row.Scan(&hash); err {
		case sql.ErrNoRows:
			fmt.Println("User doesnt exist in database")
		case nil:
			fmt.Println(hash)
		default:
			panic(err)
		}

		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
		if err != nil {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}

		http.SetCookie(w, c)

		insertQuery := `INSERT INTO sessions (value,username) VALUES($1,$2)`
		_, err = database.Conn.Exec(insertQuery, c.Value, name)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "login.html", nil)

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	c, _ := r.Cookie("session")
	// delete the session

	selectQuery := `DELETE FROM sessions WHERE value=$1;`
	_ = database.Conn.QueryRow(selectQuery, c.Value)

	// remove the cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
