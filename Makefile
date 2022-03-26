
BQ_VERSION=https://github.com/CamposBruno/browserquest
#BQ_VERSION=https://github.com/CylonicRaider/LeetQuest

basic: dep
	go build

run: basic
	./dungeonQuest -i2p=true -client=./conf/BrowserQuest/

conf/BrowserQuest:
	git clone $(BQ_VERSION) conf/BrowserQuest

bq:
	cd conf/BrowserQuest && \
		git pull && \
		npm install #&& \
		#npm run build

mapset:
	cd conf/BrowserQuest #&& \
		#npm run build-maps

conf/config.json:
	cp -v config.json conf/config.json

conf/maps:
	cp -vr maps conf/maps

conf/BrowserQuest/client/config/config_local.json: bq mapset conf/config.json conf/maps
	#mkdir -p conf/BrowserQuest/client/config/
	#cp conf/BrowserQuest/config.default.json conf/BrowserQuest/client/config/config_local.json
	#cp conf/BrowserQuest/client/config/config_local.json conf/BrowserQuest/config.json
	#cp conf/BrowserQuest/client/config/config_build.json-dist conf/BrowserQuest/client/config/config_local.json


VERSION=0.0.01
CGO_ENABLED=0
export CGO_ENABLED=0

GOOS?=$(shell uname -s | tr A-Z a-z)
GOARCH?="amd64"

ARG=-v -tags netgo -ldflags '-w -extldflags "-static"'

BINARY=dungeonQuest
SIGNER=hankhill19580@gmail.com
CONSOLEPOSTNAME=MMORPG
USER_GH=eyedeekay

build: dep
	go build $(ARG) -tags="netgo" -o $(BINARY)-$(GOOS)-$(GOARCH)
	make su3

clean:
	rm -f $(BINARY)-plugin plugin $(BINARY)-*zip -r
	rm -f *.su3 *.zip $(BINARY)-$(GOOS)-$(GOARCH) $(BINARY)-*
	git clean -xdf

distclean: clean
	rm -rfv conf/BrowserQuest

all: windows linux osx bsd

windows:
	GOOS=windows GOARCH=amd64 make build su3
	GOOS=windows GOARCH=386 make build su3

linux:
	GOOS=linux GOARCH=amd64 make build su3
	GOOS=linux GOARCH=arm64 make build su3
	GOOS=linux GOARCH=386 make build su3

osx:
	GOOS=darwin GOARCH=amd64 make build su3
	GOOS=darwin GOARCH=arm64 make build su3

bsd:
	GOOS=freebsd GOARCH=amd64 make build su3
	GOOS=openbsd GOARCH=amd64 make build su3

dep: conf/BrowserQuest conf/BrowserQuest/client/config/config_local.json
	mkdir -p conf/lib/
#	cp "$(HOME)/build/shellservice.jar" conf/lib/shellservice.jar -v

su3:
	i2p.plugin.native -name=$(BINARY)-$(GOOS)-$(GOARCH) \
		-signer=$(SIGNER) \
		-version "$(VERSION)" \
		-author=$(SIGNER) \
		-autostart=true \
		-clientname=$(BINARY)-$(GOOS)-$(GOARCH) \
		-consolename="$(BINARY) - $(CONSOLEPOSTNAME)" \
		-consoleurl="http://127.0.0.1:7681/index.html" \
		-name="$(BINARY)-$(GOOS)-$(GOARCH)" \
		-delaystart="1" \
		-desc="`cat desc`" \
		-exename=$(BINARY)-$(GOOS)-$(GOARCH) \
		-icondata=icon/icon.png \
		-updateurl="http://idk.i2p/$(BINARY)/$(BINARY)-$(GOOS)-$(GOARCH).su3" \
		-website="http://idk.i2p/$(BINARY)/" \
		-command="$(BINARY)-$(GOOS)-$(GOARCH) -i2p=true" \
		-license=MPL \
		-res=conf/
	unzip -o $(BINARY)-$(GOOS)-$(GOARCH).zip -d $(BINARY)-$(GOOS)-$(GOARCH)-zip

sum:
	sha256sum $(BINARY)-$(GOOS)-$(GOARCH).su3

version:
	gothub release -u eyedeekay -r $(BINARY) -t "$(VERSION)" -d "`cat desc`"; true

upload:
	gothub upload -R -u eyedeekay -r $(BINARY) -t "$(VERSION)" -f $(BINARY)-$(GOOS)-$(GOARCH).su3 -n $(BINARY)-$(GOOS)-$(GOARCH).su3 -l "`sha256sum $(BINARY)-$(GOOS)-$(GOARCH).su3`"

upload-windows:
	GOOS=windows GOARCH=amd64 make upload
	GOOS=windows GOARCH=386 make upload

upload-linux:
	GOOS=linux GOARCH=amd64 make upload
	GOOS=linux GOARCH=arm64 make upload
	GOOS=linux GOARCH=386 make upload

upload-osx:
	GOOS=darwin GOARCH=amd64 make upload
	GOOS=darwin GOARCH=arm64 make upload

upload-bsd:
	GOOS=freebsd GOARCH=amd64 make upload
	GOOS=openbsd GOARCH=amd64 make upload

upload-all: upload-windows upload-linux upload-osx upload-bsd

download-su3s:
	GOOS=windows GOARCH=amd64 make download-single-su3
	GOOS=windows GOARCH=386 make download-single-su3
	GOOS=linux GOARCH=amd64 make download-single-su3
	GOOS=linux GOARCH=arm64 make download-single-su3
	GOOS=linux GOARCH=386 make download-single-su3
	GOOS=darwin GOARCH=amd64 make download-single-su3
	GOOS=darwin GOARCH=arm64 make download-single-su3
	GOOS=freebsd GOARCH=amd64 make download-single-su3
	GOOS=openbsd GOARCH=amd64 make download-single-su3

download-single-su3:
	wget -N -c "https://github.com/$(USER_GH)/$(BINARY)/releases/download/$(VERSION)/$(BINARY)-$(GOOS)-$(GOARCH).su3"

release: all version upload-all

index:
	@echo "<!DOCTYPE html>" > index.html
	@echo "<html>" >> index.html
	@echo "<head>" >> index.html
	@echo "  <title>$(BINARY) - $(CONSOLEPOSTNAME)</title>" >> index.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/style.css\" />" >> index.html
	@echo "</head>" >> index.html
	@echo "<body>" >> index.html
	sed 's|goBrowserQuest|dungeonQuest|g' README.md | sed 's|SineYuan|eyedeekay|g' | pandoc  >> index.html
	@echo "</body>" >> index.html
	@echo "</html>" >> index.html

docker: dep
	docker build -t dungeonquest .

docker-run:
	docker run --restart=always -d --name dungeonQuest --net=host -v $(PWD)/conf/:/home/dungeonQuest dungeonquest