package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/michael-wang/auth"
	"github.com/michael-wang/auth/providers"
	"github.com/michael-wang/envar"
)

const (
	envPort               = "PORT"
	envGithubClientID     = "GITHUB_CLIENT_ID"
	envGithubClientSecret = "GITHUB_CLIENT_SECRET"
	envSessionSecret      = "SESSION_SECRET"
)

var (
	session *sessions.Session
)

func init() {
	envar.Load()
	envar.SetDef(envPort, 8080)
	envar.SetDef(envSessionSecret, "YYIXn86jk3GiOGzP6gnI5G8VSZaDg+qo")
}

func main() {
	flagPort := flag.Int("port", 8080, "port server listen on")
	flag.Parse()

	secret := envar.String(envSessionSecret)
	session := sessions.New([]byte(secret))
	session.Lifetime = 12 * time.Hour

	github, err := providers.Github(
		envGithubClientID,
		envGithubClientSecret,
		fmt.Sprintf("http://localhost:%d/callback/github", envar.Int(envPort)),
	)
	if err != nil {
		panic(err)
	}
	auth.AddOAuth2Provider(
		github,
		session,
		fmt.Sprintf("http://localhost:%d/", envar.Int(envPort)),
	)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(homeHandler))
	mux.Handle("/login/github", auth.LoginGithubHandler())
	mux.Handle("/callback/github", auth.CallbackHandler())

	addr := fmt.Sprintf(":%d", *flagPort)
	fmt.Println("http server listen on", addr)
	err = http.ListenAndServe(addr, logger(session.Enable(mux)))
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
