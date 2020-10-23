package main

import "github.com/zmb3/spotify"

// SpotifyClient is a wrapper interface around spotify.client
// used in order to improve testability of the code.
type SpotifyClient interface {
	Player
	Pause() error
	Previous() error
	Next() error
	PlayerCurrentlyPlaying() (*spotify.CurrentlyPlaying, error)
	PlayerDevices() ([]spotify.PlayerDevice, error)
	TransferPlayback(spotify.ID, bool) error
}

// Player 音楽を再生するやつ
type Player interface {
	Play() error
	PlayOpt(opt *spotify.PlayOptions) error
}
