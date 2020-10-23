# Gotify

Gotify is a Spotify player for cli.

Inspired by https://github.com/jbszczepaniak/spotify-cli.

Note: なんか知らんけど wsl2 環境で動作が異常(まだ原因特定できてない)

## Feature

- play
- pause
- next
- previous
- shuffle (set shuffle)

## Usage

### Require
1. Premium Spotify Account
1. Created Spotify Application under https://developer.spotify.com/dashboard/applications (set redirect URI to http://localhost:7777/gotify)


You'll get a client ID and secret key for your application. 

### Installation

```sh
$ go get github.com/zztkm/gotify
```
settings.json must be located at `~/.config/gotify/settings.json`.  
For windows `%APPDATA%\gotify/settings.json`

```json
{
    "clientID": "sssafdafdfsaf",
    "secretKey": "dafsafasf"
}
```

`gotify`