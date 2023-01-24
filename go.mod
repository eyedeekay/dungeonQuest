module github.com/eyedeekay/dungeonQuest

go 1.17

require (
	fyne.io/systray v1.10.0
	github.com/SineYuan/goBrowserQuest v0.0.0-20181227111843-ab6c4a824b1a
	github.com/atotto/clipboard v0.1.4
	github.com/bitly/go-simplejson v0.5.0
	github.com/eyedeekay/go-i2pbrowser v0.0.6
	github.com/eyedeekay/magnetWare v0.0.0-20230123050831-1461312b326e
	github.com/eyedeekay/onramp v0.0.0-20230118065332-eb11a4ec6434
	github.com/eyedeekay/unembed v0.0.0-20230123014222-9916b121855b
	github.com/gorilla/websocket v1.5.0
	github.com/labstack/echo v3.3.10+incompatible
)

require (
	github.com/artdarek/go-unzip v1.0.0 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/cretz/bine v0.2.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/eyedeekay/go-fpw v0.0.6 // indirect
	github.com/eyedeekay/i2pkeys v0.33.0 // indirect
	github.com/eyedeekay/sam3 v0.33.5 // indirect
	github.com/godbus/dbus/v5 v5.0.4 // indirect
	github.com/google/go-github v17.0.0+incompatible // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/tevino/abool v1.2.0 // indirect
	github.com/urfave/cli/v2 v2.23.7 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xgfone/bt v0.4.3 // indirect
	github.com/xgfone/bttools v0.0.0-00010101000000-000000000000 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
)

replace github.com/SineYuan/goBrowserQuest => ./

replace github.com/xgfone/bttools => ../../xgfone/bttools

replace github.com/xgfone/bt => ../../xgfone/bt

replace github.com/artdarek/go-unzip v1.0.0 => github.com/eyedeekay/go-unzip v0.0.0-20230124015700-cc3131fd4ee0

replace github.com/eyedeekay/go-i2pbrowser => ../../eyedeekay/go-i2pbrowser