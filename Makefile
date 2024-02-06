release_binary:
	go build -ldflags "-s -w"
build:
	go build
run:
	go run main.go
watch:
	air
