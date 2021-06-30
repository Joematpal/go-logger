.PHONY: example

example:
	LOG_ENV=prod \
	LOG_ENCODING=console \
	LOG_LEVEL=debug \
	go run example/main.go