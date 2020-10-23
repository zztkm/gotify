package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

const redirectURI = "http://localhost:8080/callback"

// We'll want these variables sooner rather than later
var (
	client      *spotify.Client
	playerState *spotify.PlayerState
)

var html = `
<br/>
<a href="/player/play">Play</a><br/>
<a href="/player/pause">Pause</a><br/>
<a href="/player/next">Next track</a><br/>
<a href="/player/previous">Previous Track</a><br/>
<a href="/player/shuffle">Shuffle</a><br/>

`

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

// LivePrefixState Prefix
var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

type credential struct {
	ClientID  string `json:"clientID"`
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
		os.Exit(1)
	}
	var c credential
	err = json.Unmarshal(b, &c)
	if err != nil {
		fmt.Println("Failed to unmarshal file: ", err)
		os.Exit(1)
	}
	return c
}

func executor(in string) {
	in = strings.TrimSpace(in)
	var err error
	switch in {
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	case "play":
		err = client.Play()
		LivePrefixState.LivePrefix = in + "ing > "
		LivePrefixState.IsEnable = true
	case "pause":
		err = client.Pause()
		LivePrefixState.LivePrefix = in + "d > "
		LivePrefixState.IsEnable = true
	case "next":
		err = client.Next()
	case "previous":
		err = client.Previous()
	case "shuffle":
		playerState.ShuffleState = !playerState.ShuffleState
		err = client.Shuffle(playerState.ShuffleState)
	default:
		LivePrefixState.IsEnable = false
		LivePrefixState.LivePrefix = in
		return
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "play", Description: "曲を再生するぜ！"},
		{Text: "pause", Description: "曲を一時停止するぜ！"},
		{Text: "next", Description: "次の曲を再生するぜ！"},
		{Text: "previous", Description: "前の曲を再生するぜ！"},
		{Text: "shuffle", Description: "曲をシャッフル再生するぜ！"},
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func changeLivePrefix() (string, bool) {
	return LivePrefixState.LivePrefix, LivePrefixState.IsEnable

}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
		os.Exit(1)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
		os.Exit(1)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Login Completed!"+html)
	fmt.Println("Login Completed!")
	ch <- &client
}

func main() {
	flag.Parse()

	if *printVersion {
		fmt.Printf("%s %s (rev: %s/%s)\n", name, version, revision, runtime.Version())
		return
	}

	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	c := getCredentials()
	auth.SetAuthInfo(c.ClientID, c.SecretKey)
	fmt.Println(c.ClientID, ":", c.SecretKey)
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client = <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println("You are logged in as:", user.ID)

	playerState, err = client.PlayerState()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("Found your %s (%s)\n", playerState.Device.Type, playerState.Device.Name)

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
