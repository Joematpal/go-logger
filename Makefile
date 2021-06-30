.PHONY: example

example:
	LOG_ENV=prod \
	LOG_ENCODING=json \
	LOG_LEVEL=debug \
	LST=false \
	go run example/main.go