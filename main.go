package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/zmb3/spotify"
)

const (
	name     = "spotify for cli"
	version  = "0.1.0"
	revision = "HEAD"
)

var playerState *spotify.PlayerState

// LivePrefixState Prefix
var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

type credential struct {
	ClientID  string `json:"clinetID"`
	SecretKey string `json:"secretKey"`
}

var printVersion = flag.Bool("version", false, "print version")

func createConfigFileName() string {
	file := "settings.json"

	if runtime.GOOS == "windows" {
		file = filepath.Join(os.Getenv("APPDATA"), "spotify", file)
	} else {
		file = filepath.Join(os.Getenv("HOME"), ".config", "spotify", file)
	}

	return file
}

func getCredentials() credential {
	file := createConfigFileName()

	b, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Failed to read config file: ", err)
	}
	var c credential
	err = json.Unmarshal(b, &c)
	if err != nil {
		fmt.Println("Failed to unmarshal file: %s\n", err)
	}
	return c
}

func newSpotifyAuthenticator() spotify.Authenticator {
	redirectURI := url.URL{Scheme: "http", Host: "localhost:8888", Path: "/spotify"}

	auth := spotify.NewAuthenticator(
		redirectURI.String(),
		spotify.ScopeUserReadPrivate,
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
		spotify.ScopeUserModifyPlaybackState,
		spotify.ScopeUserLibraryRead,
		// Used for Web Playback SDK
		"streaming",
		spotify.ScopeUserReadEmail,
	)
	credentials := getCredentials()
	auth.SetAuthInfo(credentials.ClientID, credentials.SecretKey)
	return auth
}

var client SpotifyClient
var spotifyAuthenticator = newSpotifyAuthenticator()

func executor(in string) {
	in = strings.TrimSpace(in)

	switch in {
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	case "play":

	default:
		LivePrefixState.IsEnable = false
		LivePrefixState.LivePrefix = in
		return
	}

	LivePrefixState.LivePrefix = in + "> "
	LivePrefixState.IsEnable = true
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
		{Text: "groups", Description: "Combine users with specific rules"},
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func changeLivePrefix() (string, bool) {
	return LivePrefixState.LivePrefix, LivePrefixState.IsEnable

}

func main() {
	flag.Parse()

	if *printVersion {
		fmt.Printf("%s %s (rev: %s/%s)\n", name, version, revision, runtime.Version())
		return
	}
	var client SpotifyClient
	var spotifyAuthenticator = newSpotifyAuthenticator()

	authHandler := &web.AuthHandler{
		Client:        make(chan *spotify.Client),
		State:         uuid.New().String(),
		Authenticator: spotifyAuthenticator,
	}

	webSocketHandler := &web.WebsocketHandler{
		PlayerShutdown:    make(chan bool),
		PlayerDeviceID:    make(chan spotify.ID),
		PlayerStateChange: make(chan *web.WebPlaybackState),
	}

	if debugMode {
		client = player.NewDebugClient()
		go func() {
			webSocketHandler.PlayerDeviceID <- "debug"
		}()
	} else {
		var err error

		h := http.NewServeMux()
		h.Handle("/ws", webSocketHandler)
		h.Handle("/spotify-cli", authHandler)
		h.HandleFunc("/player", web.PlayerHandleFunc)

		go func() {
			log.Fatal(http.ListenAndServe(":8888", h))
		}()

		err = player.StartRemoteAuthentication(spotifyAuthenticator, authHandler.State)
		if err != nil {
			log.Printf("could not get client, shutting down, err: %v", err)
		}
	}

	// wait for authentication to complete
	client = <-authHandler.Client

	// wait for device to be ready
	webPlayerID := <-webSocketHandler.PlayerDeviceID

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("> "),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionPrefixTextColor(prompt.Yellow), // Prefix(ここでは >>>) の色を黄色に変更
		prompt.OptionTitle("spotify cli"),
	)
	p.Run()
}
