# matterpoll-emoji

Polling feature for mattermost's custom slash command.

## Setup server

```
git clone https://github.com/kaakaa/matterpoll-emoji.git
cd matterpoll-emoji
glide install
go run main.go -p 8066
```

## Setup mattermost

Create a `Custome Slash Command` from Integration > Slash Commands > Add Slash Command.

* DisplayName - Arbitrary (ex. MatterPoll)
* Description - Arbitrary (ex. Polling feature by https://github.com/kaakaa/matterpoll-emoji)
* Command Trigger Word - `poll`
* Request URL - http://localhost:8066
* Request Method - `POST`
* Others - optional

## Usage

Typing this on mattermost

```
/poll `What do you gys wanna grab for lunch?` :pizza: :sushi: :fried_shrimp: :spaghetti: :apple:
```

then posting poll comment

![screen_shot](https://raw.githubusercontent.com/kaakaa/matterpoll-emoji/master/matterpoll-emoji.png)

## License
* MIT
  * see [LICENSE](LICENSE)

