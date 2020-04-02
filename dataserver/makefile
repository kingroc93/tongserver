# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME=tongds
BINARY_UNIX=$(BINARY_NAME)_linux

# outPath
OUT=out

clean:
	rm out/* -rf

build: clean
	$(GOBUILD) -o $(OUT)/$(BINARY_NAME) -v
	mkdir $(OUT)/conf
	cp conf/* $(OUT)/conf -r

run: build
	cd $(OUT) && ./$(BINARY_NAME)

crs: clean
	mkdir $(OUT)/conf
	cp conf/* $(OUT)/conf -r
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -ldflags '-extldflags "-static"' -installsuffix cgo -o $(OUT)/$(BINARY_UNIX) -v
	sudo docker build -t golang-scratch-app .

