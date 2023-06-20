APPNAME="Port Punch"
APPID="com.nbembedded.portpunch"
VERSION?=development
BUILD_FLAGS="-X=main.version=$(VERSION)"
OUT_DIR?=bin

all: build

build:
	go build $(BUILD_FLAGS) -o ./$(OUT_DIR)/

package:
	~/go/bin/fyne-cross darwin -icon icon.png --app-id $(APPID) --name $(APPNAME) -ldflags=$(BUILD_FLAGS)
	~/go/bin/fyne-cross darwin -icon icon.png --app-id $(APPID) --name $(APPNAME) -arch=amd64 -ldflags=$(BUILD_FLAGS)
	~/go/bin/fyne-cross windows --app-id $(APPID) --name $(APPNAME) -arch=amd64 -ldflags=$(BUILD_FLAGS)
test:
	go test ./... -v

clean:
	rm -f $(OUT_DIR)/*

.PHONY: all build package test clean
