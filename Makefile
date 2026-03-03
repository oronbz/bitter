.PHONY: build run clean

build:
	go build -o bitter .

run:
	go run .

clean:
	rm -f bitter
