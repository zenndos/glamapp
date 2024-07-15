.PHONY: build run dev clean

build:
	go build -o glamapp

run: build
	./glamapp

dev:
	go run main.go

clean:
	rm -f glamapp