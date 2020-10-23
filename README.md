# Spotify player

SpotifyをCLIで操作したくない？

## Installation

```sh
$ go get github.com/zztkm/go-spotify
```

## Usage

まずはセッティング！

以下のページにアクセスして、自分のアプリケーション(ex. MySpotifyPlayer)

https://developer.spotify.com/my-applications/.

You'll get a client ID and secret key for your application. 

client ID と secret key を下記のように`settings.json`に書き込みます。

UNIX:
```
~/.config/spotify/settings.json
```

WINDOWSのコンフィグ配置場所: 
```
%APPDATA%\spotify/settings.json
```

セッティング例: 
```json
{
    "clientID": "sssafdafdfsaf",
    "secretKey": "dafsafasf"
}
```

あとは実行！`spotify`