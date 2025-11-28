.PHONY: format run-local -test

format:
	go fmt

run-local:
	go run .

demo:
	go run . -demo

test:
	echo TODO
