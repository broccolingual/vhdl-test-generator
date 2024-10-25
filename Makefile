.PHONY: clean install run build

clean:
	rm -rf gen.exe

install:
	# install 3rd party library
	# go get ~

run:
	go run main.go

build: install clean
	go build -ldflags="-s -w -buildid=" -trimpath -o gen.exe main.go
