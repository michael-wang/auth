package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/michael-wang/auth"
)

const (
	// TODO get following values from env var
	githubClientID     = "faf9d0195aed6907ee12"
	githubClientSecret = "2e1e0b8b51d9dcb8c8ef9c7def16383632f20b23"
	sessionSecret      = "YYIXn86jk3GiOGzP6gnI5G8VSZaDg+qo"
)

var (
	session *sessions.Session
)

func main() {
	flagPort := flag.Int("port", 8080, "port server listen on")
	flag.Parse()

	session := sessions.New([]byte(sessionSecret))
	session.Lifetime = 12 * time.Hour

	// TODO use 'port' flag for redirect and after login URL.
	a := auth.New(session, githubClientID, githubClientSecret, "http://localhost:8080/callback/github", "http://localhost:8080/")

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(homeHandler))
	mux.Handle("/login/github", a.LoginGithubHandler())
	mux.Handle("/callback/github", a.OAuth2CallbackHandler())

	addr := fmt.Sprintf(":%d", *flagPort)
	fmt.Println("http server listen on", addr)
	err := http.ListenAndServe(addr, logger(session.Enable(mux)))
	if err != nil {
		panic(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	bb := session.GetBytes(r, auth.LoginUser)
	if bb == nil {
		// not logged in
		// TODO get url based on 'port' flag.
		fmt.Fprint(w, `<a href="http://localhost:8080/login/github">Login Github</a>`)
		return
	}

	user := &auth.GithubUser{}
	json.Unmarshal(bb, user)
	fmt.Fprintf(w, "Hi %s (%s)\n", user.Name, user.Email)
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s\t%s\t%s\n", r.Proto, r.Method, r.URL)

		next.ServeHTTP(w, r)
	})
}
