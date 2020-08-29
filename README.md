# gamejam scaffold (WIP)

## Server
- [centrifugo](https://github.com/centrifugal/centrifugo)

## Client game libary
- [ebiten](https://github.com/hajimehoshi/ebiten)

## Usage
### Server
#### Install
```
git clone https://github.com/centrifugal/centrifugo
cd centrifugo
go build .
./centrifugo genconfig
```
#### Config Setting
add some setting to config.json
```
{
  ....
  "admin": true,
  "publish": true
}
```
#### Run
```
./centrifugo --config=config.json
```
[server admin](http://localhost:8000)

### Client
```
git clone https://github.com/miles990/gamejam.git
cd gamejam
## Change your exampleTokenHmacSecret variable in /game/event_handler.go
## It's the same setting in server config.json 
## "token_hmac_secret_key"
go run main.go
```