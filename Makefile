.PHONY: build run dev clean

build:
	go build -o glamapp

run: build
	./glamapp

dev:
	go run src/main.go

clean:
	rm -f glamapp