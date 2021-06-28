package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

type Auth struct {
	config *oauth2.Config
}

func New(clientID, clientSecret, callbackURL string) *Auth {
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
	}
}

func (a *Auth) GetLoginURL() string {
	return a.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

type UserProfile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (a *Auth) GetUserProfile(code string) (*UserProfile, error) {
	ctx := context.Background()
	token, err := a.config.Exchange(ctx, code)
	if err != nil {
		panic(err)
	}

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

	user := &UserProfile{}
	err = json.Unmarshal(bb, user)
	if err != nil {
		panic(err)
	}
	return user, nil
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.Form["code"][0]
	fmt.Fprintf(w, "Code: %s\n\n", code)
	fmt.Fprint(w, "Select and Copy code above, and press to testing program")
}
