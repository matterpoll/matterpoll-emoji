# matterpoll-emoji

[![Build Status](https://travis-ci.org/kaakaa/matterpoll-emoji.svg?branch=master)](https://travis-ci.org/kaakaa/matterpoll-emoji)
[![Code Coverage](https://codecov.io/gh/kaakaa/matterpoll-emoji/branch/master/graph/badge.svg)](https://codecov.io/gh/kaakaa/matterpoll-emoji/branch/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaakaa/matterpoll-emoji)](https://goreportcard.com/report/github.com/kaakaa/matterpoll-emoji)
[![Releases](https://img.shields.io/github/release/kaakaa/matterpoll-emoji.svg)](https://github.com/kaakaa/matterpoll-emoji/releases/latest)

Polling feature for Mattermost's custom slash command.

## Requirements
- [go](https://golang.org/)
    - Verion 1.7 or newer
- [Mattermost server](https://about.mattermost.com/)
    - Version 3.7 or newer

## Setup Guide

### Setup Mattermost

Create a `Custom Slash Command` from Integration > Slash Commands > Add Slash Command.

* DisplayName - `Matterpoll`
* Description - `Polling feature by https://github.com/kaakaa/matterpoll-emoji`
* Command Trigger Word - `poll`
* Request URL - `http://localhost:8505/poll`
* Request Method - `POST`
* Response Username - Leave this empty
* Response Icon - Leave this empty
* Autocomplete - Enable this
* Autocomplete Hint - `[Question] [Option1] [Option2]...`
* Autocomplete Description - `Start a poll`

Copy the Token from your newly created slash command

**Caution**: If you run `matterpoll-emoji` on same host as Mattermost server, you have to add `localhost` to [**Allow untrusted internal connections to**](https://docs.mattermost.com/administration/config-settings.html#allow-untrusted-internal-connections-to) option.

### Setup server

#### Run pre compiled release

Download the latest version at https://github.com/kaakaa/matterpoll-emoji/releases/latest.
Decompress it and change parameter in `config.json` as you need them
```
{
  "host": "http://mattermost.example.com:8065",  // The URL of your Mattermost server
  "listen": "localhost:8505",  // The address:port to listen on
  "token": "9jrxak1ykxrmnaed9cps9i4cim",  // The Token created my Mattermost
  "user": {
   "id": "bot",          // The username of an existing Mattermost account
   "password": "botbot"  // The password of an existing Mattermost account
 }
}
```
Run the server
```
./matterpoll-emoji
```

#### Compile the source my yourself

Clone this repository and checkout the latest release. You can just use the master branch but it can be unstable.
```
go get -u https://github.com/kaakaa/matterpoll-emoji
cd $GOPATH/src/github.com/kaakaa/matterpoll-emoji
git checkout $(git describe --tags)
```
Copy the default config
```
cp .config.json config.json
```
Change parameter in `config.json` as you need them
```
{
  "host": "http://mattermost.example.com:8065",  // The URL of your Mattermost server
  "listen": "localhost:8505",  // The address:port to listen on
  "token": "9jrxak1ykxrmnaed9cps9i4cim",  // The Token created my Mattermost
  "user": {
   "id": "bot",          // The username of an existing Mattermost account
   "password": "botbot"  // The password of an existing Mattermost account
 }
}
```

Run the server
```
make run
```

## Usage

Typing this on Mattermost

```
/poll `What do you gys wanna grab for lunch?` :pizza: :sushi: :fried_shrimp: :spaghetti: :apple:
```

then posting poll comment

![screen_shot](https://raw.githubusercontent.com/kaakaa/matterpoll-emoji/master/matterpoll-emoji.png)

You can use `"` or `'` instead of `` ` ``

## License
* MIT
  * see [LICENSE](LICENSE)
