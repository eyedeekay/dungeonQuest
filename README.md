dungeonQuest
============

go implementation for [BrowserQuest](https://github.com/mozilla/BrowserQuest) server,
with special features for use with I2P.

Installation
-------------

```
go get github.com/eyedeekay/dungeonQuest
```

Configuration
-------------

```
Usage of ./dungeonQuest:
  -client string
    	BrowserQuest root directory to serve if provided (default "./BrowserQuest")
  -config string
    	configuration file path (default "./config.json")
  -i2p
    	use I2P
  -port string
    	port to present the plugin homepage on, actually a link to the game. (default "7681")
  -prefix string
    	request url prefix when client is provided, cannot be '/'  (default "/game")
  -tls
    	use TLS (default true)
```

Deployment
----------

### client 
```
git clone https://github.com/mozilla/BrowserQuest.git

## or obtain an updated fork here:

git clone https://github.com/CamposBruno/browserquest

cp BrowserQuest/client/config/config_local.json-dist BrowserQuest/client/config/config_local.json 
```
edit `BrowserQuest/client/config/config_local.json` to set server host and port.

### server

```
cd $GOPATH/src/github.com/SineYuan/goBrowserQuest
go build main.go
./main -config /path/to/config.json -client /path/to/BrowserQuest 
```

### docker

```
docker build -t gobrowserquest .
docker run --restart=always -d --name dungeonQuest -v $(PWD)/conf/:/home/dungeonQuest gobrowserquest
```

then you can play game at `http://{HOST}:{PORT}/game/client/index.html`

### I2P Plugin

- [http://idk.i2p/dungeonQuest/dungeonQuest-darwin-amd64.su3](http://idk.i2p/dungeonQuest/dungeonQuest-darwin-amd64.su3)
- [http://idk.i2p/dungeonQuest/dungeonQuest-darwin-arm64.su3](http://idk.i2p/dungeonQuest/dungeonQuest-darwin-arm64-.su3)
- [http://idk.i2p/dungeonQuest/dungeonQuest-linux-amd64.su3](http://idk.i2p/dungeonQuest/dungeonQuest-linux-amd64.su3)
- [http://idk.i2p/dungeonQuest/dungeonQuest-linux-arm64.su3](http://idk.i2p/dungeonQuest/dungeonQuest-linux-arm64.su3)
- [http://idk.i2p/dungeonQuest/dungeonQuest-linux-386.su3](http://idk.i2p/dungeonQuest/dungeonQuest-linux-386.su3)
- [http://idk.i2p/dungeonQuest/dungeonQuest-freebsd-amd64.su3](http://idk.i2p/dungeonQuest/dungeonQuest-freebsd-amd64.su3)
- [http://idk.i2p/dungeonQuest/dungeonQuest-openbsd-amd64.su3](http://idk.i2p/dungeonQuest/dungeonQuest-openbsd-amd64-.su3)
- [http://idk.i2p/dungeonQuest/dungeonQuest-windows-amd64.su3](http://idk.i2p/dungeonQuest/dungeonQuest-windows-amd64.su3)
- [http://idk.i2p/dungeonQuest/dungeonQuest-windows-386.su3](http://idk.i2p/dungeonQuest/dungeonQuest-windows-386.su3)

### Docker Container with I2P

```
docker build -t gobrowserquest .
docker run --restart=always --net=host -d --name dungeonQuest -v $(PWD)/conf/:/home/dungeonQuest gobrowserquest dungeonQuest -i2p -client /home/dungeonQuest/BrowserQuest -config /home/dungeonQuest/config.json
```

```
then you can play game at `http://{BASE32}:8000/game/client/index.html`
```

TODO
----------
goBrowserQuest have yet to implement all the function of BrowserQuest server. welcome to forks