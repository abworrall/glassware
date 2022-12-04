package spot // Named spot to avoid collisions with library

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	ctx := context.Background()

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
	client := spotify.New(auth.Client(ctx, t))

	// Did we just refresh ? If so, save the new token
	t2, err := client.Token()
	if err == nil && t.Expiry != t2.Expiry {
		log.Printf("    Token was refreshed under the hood, saving\n")
		saveToken(c.CacheDir, t2)
	}

	if user, err := client.CurrentUser(ctx); err == nil {
		log.Printf("    Have a spotify client logged in as: %s\n", user.DisplayName)
	} else {
		log.Printf("Spotify client broken with CurrentUser: %s\n\nToken: %#v\n\nClient: %#v\n\n", err, t, client)
		return nil
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

	log.Printf("    Stored the OAuth2 token: %s (expires in %s, at %s)\n", filename, time.Until(t.Expiry), t.Expiry)
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

	log.Printf("    Loaded OAuth2 token: %s (expires in %s, at %s)\n", filename, time.Until(t.Expiry), t.Expiry)

	return &t
}

func StoreCreds(cacheDir, id, secret string) error {
	saveFile := func(filename, contents string) error {
		filepath := cacheDir + "/" + filename
		if err := os.WriteFile(filepath, []byte(contents), 0600); err != nil {
			return err
		}
		return nil
	}

	if err := saveFile("id", id); err != nil {
		return err
	}
	if err := saveFile("secret", secret); err != nil {
		return err
	}

	return nil
}

func LoadCreds(cacheDir string) (string, string, error) {
	loadFile := func(filename string) (string, error) {
		filepath := cacheDir + "/" + filename
		if b, err := os.ReadFile(filepath); err != nil {
			return "", err
		} else {
			return string(b), nil
		}
	}

	id, err1 := loadFile("id")
	secret, err2 := loadFile("secret")
	if err1 != nil { return "", "", err1 }
	if err2 != nil { return "", "", err2 }

	return id, secret, nil
}
