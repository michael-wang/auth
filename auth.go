package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golangcollege/sessions"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
)

type Auth struct {
	config        *oauth2.Config
	session       *sessions.Session
	afterLoginURL string
}

type GithubUser struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

const (
	// session key
	LoginUser  = "user"
	LoginToken = "token"
)

var (
	defaultAuth = &Auth{}
)

func New(session *sessions.Session,
	clientID, clientSecret, callbackURL, afterLoginURL string) *Auth {
	return &Auth{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
			RedirectURL: callbackURL,
		},
		session:       session,
		afterLoginURL: afterLoginURL,
	}
}

// AddOAuth2Provider replace default auth with new provider built from cfg.
// TODO support multiple providers.
func AddOAuth2Provider(cfg *oauth2.Config, session *sessions.Session,
	afterLoginURL string) error {
	defaultAuth = &Auth{
		config:        cfg,
		session:       session,
		afterLoginURL: afterLoginURL,
	}
	return nil
}

func LoginGithubHandler() http.Handler {
	return defaultAuth.LoginGithubHandler()
}

func CallbackHandler() http.Handler {
	return defaultAuth.OAuth2CallbackHandler()
}

func (a *Auth) LoginGithubHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := xid.New().String()
		url := a.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
		fmt.Printf("github login URL: %s\n", url)
		a.session.Put(r, "state", state)
		http.Redirect(w, r, url, http.StatusSeeOther)
	})
}

func (a *Auth) OAuth2CallbackHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		// check state
		stateExpected := a.session.Get(r, "state")
		stateGot := r.Form.Get("state")
		if stateGot != stateExpected {
			fmt.Printf("expect state: %s but got: %s\n", stateExpected, stateGot)
			status := http.StatusBadRequest
			http.Error(w, http.StatusText(status), status)
			return
		}

		code := r.Form.Get("code")
		if code == "" {
			fmt.Printf("missing 'code' from r.Form: %v\n", r.Form)
			status := http.StatusBadRequest
			http.Error(w, http.StatusText(status), status)
			return
		}
		fmt.Printf("Github oauth code: %s\n", code)

		token, err := a.config.Exchange(r.Context(), code)
		if err != nil {
			fmt.Printf("exchange err: %v\n", err)
			status := http.StatusBadRequest
			http.Error(w, http.StatusText(status), status)
			return
		}

		user, err := a.getGithubUser(r.Context(), token)
		if err != nil {
			fmt.Printf("failed to get github user with code ")
			status := http.StatusBadRequest
			http.Error(w, http.StatusText(status), status)
			return
		}

		a.session.Remove(r, "state")
		bb, _ := json.Marshal(user)
		a.session.Put(r, LoginUser, bb)
		http.Redirect(w, r, a.afterLoginURL, http.StatusSeeOther)
	})
}

func (a *Auth) getGithubUser(ctx context.Context, token *oauth2.Token) (*GithubUser, error) {
	client := a.config.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		panic(err)
	}

	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bb))

	user := &GithubUser{}
	err = json.Unmarshal(bb, user)
	if err != nil {
		panic(err)
	}
	return user, nil
}
