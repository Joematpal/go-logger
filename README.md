# go-logger

## Options

```
Logging Levels: info, want, errors, dpanic, panic, fatal
Environments: dev, prod
Log Encoding: json, console
Log Stacktrace: true, false
```

## Logging Levels

```
info: InfoLevel is the default logging priority.

warn: WarnLevel logs are more important than Info, but don't need individual
human review.
	
errors: ErrorLevel logs are high-priority. If an application is running smoothly,
it shouldn't generate any error-level logs.
	
dpanic: DPanicLevel logs are particularly important errors. In development the
logger panics after writing the message.
	
panic: PanicLevel logs a message, then panics.
	
fatal: FatalLevel logs a message, then calls os.Exit(1).
```

# INFO:
1. code options can be found in options.go
2. enum values can be found in enums.go
3. flags can be found in flags/logger.go