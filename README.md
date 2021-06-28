# go-logger


## Logging Levels

```
InfoLevel is the default logging priority.

WarnLevel logs are more important than Info, but don't need individual
human review.
	
ErrorLevel logs are high-priority. If an application is running smoothly,
it shouldn't generate any error-level logs.
	
DPanicLevel logs are particularly important errors. In development the
logger panics after writing the message.
	
PanicLevel logs a message, then panics.
	
FatalLevel logs a message, then calls os.Exit(1).
```