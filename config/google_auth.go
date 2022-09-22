package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleSignInRedirectURL = "http://localhost:8080/google/callback"

func SetupGoogleAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     getEnv("GOOGLE_CLIENT_ID", "client_id"),
		ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", "client_secret"),
		RedirectURL:  googleSignInRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
