# Gotify

Gotify is a Spotify player for cli.

Inspired by https://github.com/jbszczepaniak/spotify-cli.

## Feature

- play
- pause
- next
- previous
- shuffle (set shuffle)

## Usage

### require
1. Premium Spotify Account
1. Created Spotify Application under https://beta.developer.spotify.com/dashboard/applications (set redirect URI to http://localhost:8888/gotify)


You'll get a client ID and secret key for your application. 

client ID と secret key を下記のように`settings.json`に書き込みます。

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