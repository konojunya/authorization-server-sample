package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/RangelReale/osin/example"

	mysql "github.com/felipeweb/osin-mysql"
)

var (
	db *sql.DB
)

func main() {
	urlDB := "hal:hal@tcp(localhost:3306)/oauth-sample?parseTime=true"
	db, err := sql.Open("mysql", urlDB)
	if err != nil {
		panic(err)
	}

	store := mysql.New(db, "hal_")
	err = store.CreateSchemas()
	if err != nil {
		panic(err)
	}

	server := osin.NewServer(osin.NewServerConfig(), store)

	// routing
	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		resp := server.NewResponse()
		defer resp.Close()

		if ar := server.HandleAuthorizeRequest(resp, r); ar != nil {
			if !example.HandleLoginPage(ar, w, r) {
				return
			}
			ar.Authorized = true
			server.FinishAuthorizeRequest(resp, r, ar)
		}
		osin.OutputJSON(resp, w, r)
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		resp := server.NewResponse()
		defer resp.Close()

		if ar := server.HandleAccessRequest(resp, r); ar != nil {
			ar.Authorized = true
			server.FinishAccessRequest(resp, r, ar)
		}
		if resp.IsError && resp.InternalError != nil {
			log.Println("ERROR: %s", resp.InternalError)
		}
		osin.OutputJSON(resp, w, r)
	})

	fmt.Println("oauth application is listen on http://localhost:14000")
	if err := http.ListenAndServe(":14000", nil); err != nil {
		panic(err)
	}

}

func renderLoginPage(ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()

	var password string
	if err := db.QueryRow("SELECT password FROM user WHERE student_id = ?", r.Form.Get("login")).Scan(&password); err != nil {
		fmt.Println(err)
	}

	if r.Method == "POST" && r.Form.Get("password") == password {
		return true
	}

	w.Write([]byte("<html><body>"))

	w.Write([]byte(fmt.Sprintf("LOGIN %s (use test/test)<br/>", ar.Client.GetId())))
	w.Write([]byte(fmt.Sprintf("<form action=\"/authorize?%s\" method=\"POST\">", r.URL.RawQuery)))

	w.Write([]byte("Login: <input type=\"text\" name=\"login\" /><br/>"))
	w.Write([]byte("Password: <input type=\"password\" name=\"password\" /><br/>"))
	w.Write([]byte("<input type=\"submit\"/>"))

	w.Write([]byte("</form>"))

	w.Write([]byte("</body></html>"))

	return false

}
