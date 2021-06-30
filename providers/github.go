package providers

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"
)

func Github(idKey, secretKey, callbackURL string) (*oauth2.Config, error) {
	id, ok := os.LookupEnv(idKey)
	if !ok {
		return nil, fmt.Errorf("failed to find client ID by eviroment variable key: %s", idKey)
	}

	secret, ok := os.LookupEnv(secretKey)
	if !ok {
		return nil, fmt.Errorf("failed to find client secret by eviroment variable key: %s", secretKey)
	}

	return &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		RedirectURL: callbackURL,
	}, nil
}
