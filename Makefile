generate:
	go test github.com/lyft/flytectl/cmd --update

compile:
	go build -o bin/flytectl main.go
