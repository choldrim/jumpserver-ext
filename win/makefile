build:
	GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -o Jumpserver.exe main.go

build-cli:
	GOOS=windows GOARCH=amd64 go build -o Jumpserver.exe main.go

tar:
	mkdir -p Jumpserver
	cp Jumpserver.exe Jumpserver/
	cp config_example.ini Jumpserver/config.ini
	cp jms.reg Jumpserver/
	cp README.txt Jumpserver/
	zip Jumpserver.zip Jumpserver/*
	rm -rf Jumpserver
