.PHONY: example

example:
	LOG_ENV=prod \
	LOG_ENCODING=console \
	LOG_LEVEL=fatal \
	go run example/main.go