# matterpoll-emoji

[![Build Status](https://travis-ci.org/kaakaa/matterpoll-emoji.svg?branch=master)](https://travis-ci.org/kaakaa/matterpoll-emoji)
[![Code Coverage](https://codecov.io/gh/kaakaa/matterpoll-emoji/branch/master/graph/badge.svg)](https://codecov.io/gh/kaakaa/matterpoll-emoji/branch/master)

Polling feature for Mattermost's custom slash command.

## Requirements
- [go](https://golang.org/)
    - Verion 1.7 or 1.8
- [Mattermost server](https://about.mattermost.com/)
    - Version 3.0 or newer

### Setup Mattermost

Create a `Custom Slash Command` from Integration > Slash Commands > Add Slash Command.

* DisplayName - Arbitrary (ex. MatterPoll)
* Description - Arbitrary (ex. Polling feature by https://github.com/kaakaa/matterpoll-emoji)
* Command Trigger Word - `poll`
* Request URL - http://localhost:8505/poll
* Request Method - `POST`
* Others - optional

Copy the Token from your newly created slash command

### Setup server

Clone this repository
```
git clone https://github.com/kaakaa/matterpoll-emoji.git
cd matterpoll-emoji
cp .config.json config.json
```
Change parameter in `config.json` e.g.
```
{
  "host": "http://mattermost.example.com:8065", # The URL of your Mattermost server
  "token": "9jrxak1ykxrmnaed9cps9i4cim", # The Token created my Mattermost
  "user": {
   "id": "bot",          # The username of an existing Mattermost account
   "password": "botbot"  # The password of an existing Mattermost account
 }
}
```

Setup `matterpoll-emoji` server
```
glide install
go run main.go -p 8505
```

## Usage

Typing this on Mattermost

```
/poll `What do you gys wanna grab for lunch?` :pizza: :sushi: :fried_shrimp: :spaghetti: :apple:
```

then posting poll comment

![screen_shot](https://raw.githubusercontent.com/kaakaa/matterpoll-emoji/master/matterpoll-emoji.png)

## License
* MIT
  * see [LICENSE](LICENSE)
