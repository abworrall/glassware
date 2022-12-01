package spot // Named spot to avoid collisions with library

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/abworrall/glassware/pkg/config"
)

var (
	redirectPort = 8081
	redirectURIStem = "/oauth-callback"
	// Add this URI to Spotify dev portal
	redirectURI = fmt.Sprintf("http://localhost:%d%s", redirectPort, redirectURIStem)

	tokenFilename = "spotify-oauth-token.json"

	scopes = []string{
		spotifyauth.ScopePlaylistReadPrivate,
		spotifyauth.ScopePlaylistReadCollaborative,
		spotifyauth.ScopeUserReadPlaybackState,
		spotifyauth.ScopeUserModifyPlaybackState,
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopeStreaming,
	}
)

func GetClient(c config.Config, forceNewNewToken bool) *spotify.Client {
	auth := spotifyauth.New(
		spotifyauth.WithClientID(c.SpotifyId),
		spotifyauth.WithClientSecret(c.SpotifySecret),
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserModifyPlaybackState,
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeStreaming))

	// We cache the whole oauth2 token as it has a refreshtoken inside, and will work indefinitely
	t := loadToken(c.CacheDir)
	if forceNewNewToken || t == nil {
		t = doOAuthFlow(auth)
		saveToken(c.CacheDir, t)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(context.Background(), t))

	if user, err := client.CurrentUser(context.Background()); err == nil {
		log.Printf("    Have a spotify client logged in as: %s\n", user.DisplayName)
	}

	return client
}

func doOAuthFlow(auth *spotifyauth.Authenticator) *oauth2.Token {
	tokenChan := make(chan *oauth2.Token)
	state := "OauthRules"

	http.HandleFunc(redirectURIStem, func(w http.ResponseWriter, r *http.Request) {
		t, err := auth.Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}
		if st := r.FormValue("state"); st != state {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, state)
		}
		fmt.Fprintf(w, "Login Completed!")

		tokenChan <- t
	})

	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", redirectPort), nil))
	}()

	url := auth.AuthURL(state)
	log.Printf("Please log in to Spotify by visiting the following page in your browser: %s\n", url)

	// wait for auth to complete
	t := <-tokenChan
	return t
}

func saveToken(cacheDir string, t *oauth2.Token) {
	filename := cacheDir + "/" + tokenFilename

	b, err := json.Marshal(t)
	if err != nil { log.Fatal(err) }

	if err := os.WriteFile(filename, b, 0600); err != nil {
		log.Fatal(err)
	}

	log.Printf("    Stored the OAuth2 token: %s\n", filename)
}

func loadToken(cacheDir string) *oauth2.Token {
	filename := cacheDir + "/" + tokenFilename

	b, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		log.Fatal(err)
	}

	t := oauth2.Token{}
	if err := json.Unmarshal(b, &t); err != nil {
		log.Fatal(err)
	}

	log.Printf("    Loaded OAuth2 token: %s\n", filename)

	return &t
}

