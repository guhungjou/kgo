

all: kgo

kgo:
	CGO_ENABLED=0 go build -o bin/kgo
