package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

var (
	config oauth2.Config
	db     *sql.DB
)

func init() {
	urlDB := "hal:hal@tcp(localhost:3306)/oauth-sample?parseTime=true"
	conn, err := sql.Open("mysql", urlDB)
	if err != nil {
		panic(err)
	}
	db = conn

	config = oauth2.Config{
		ClientID:     "clientid",
		ClientSecret: "asdfghjkl",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:14000/authorize",
			TokenURL: "http://localhost:14000/token",
		},
		RedirectURL: "http://localhost:9000/oauth",
	}
}

func main() {

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		redirect_url := config.AuthCodeURL("hal")
		http.Redirect(w, r, redirect_url, http.StatusTemporaryRedirect)
	})
	http.HandleFunc("/oauth", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		token, err := config.Exchange(oauth2.NoContext, code)
		if err != nil {
			panic(err)
		}

		fmt.Println("access token: ", token.AccessToken)
	})

	http.ListenAndServe(":9000", nil)

}
