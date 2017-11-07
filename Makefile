build:
	go build -o mario github.com/zenja/mario/cmd;

run: build; ./mario
