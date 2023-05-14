# Personal LogStack

In the doc of zap , i couldn't easily find the way to configure. So i decide to make the wrapper around zap and zapcore package to simplify that. Also adding some other component of the log stack to make it easier for archival

The package provides a logger with different log levels (Info, Warn, Error, DPanic, Panic, Fatal, and Debug) that can be configured with different options such as enabling/disabling caller info, adding stack traces, and adding fields to the logger.

The logger instance can be used to log messages at different level :

```go
logger.Info("This is an info message")
logger.Warn("This is a warning message")
logger.Error("This is an error message")
logger.DPanic("This is a development panic message")
logger.Panic("This is a panic message")
logger.Fatal("This is a fatal message")
logger.Debug("This is a debug message")
```

## Adding Fields to the Logger

You can add fields to the logger with

```go
 logger := logStack.NewLogger(os.Stdout, logStack.InfoLevel, logStack.AddFields(
    logStack.String("key", "value"),
    logStack.Int("count", 123),
))
```

## Enabling Caller Info and Stack Trace :

```go
logger := logStack.NewLogger(os.Stdout, logStack.InfoLevel, logStack.WithCaller(true), logStack.AddStacktrack)
```

## Create Multi Logger

You can create a logger that duplicates log entries to multiple log files using the NewMultiLoger function:

```go
logs := []logStack.MultiOption{
    {
        W:     os.Stdout,
        Level: logStack.LevelEnablerFunc(func(lvl logStack.Loglevel) bool { return lvl >= logStack.DebugLevel }),
    },
    {
        W:     &lumberjack.Logger{Filename: "app.log", MaxSize: 100, MaxBackups: 10, MaxAge: 30},
        Level: logStack.LevelEnablerFunc(func(lvl logStack.Loglevel) bool { return lvl >= logStack.ErrorLevel }),
    },
}
logger := logStack.NewMultiLoger(logs)
```

You can also add rotate functionality to a multi-logger using the MultiOptionWithRotate option:

```go
logs := []logStack.MultiOptionWithRotate{
    {
        Filename: "/var/log/app.log",
        Ropt: logStack.RotateOptions{
            MaxSize:    100,
            MaxAge:     7,
            MaxBackups: 10,
            Compress:   true,
        },
        Level: logStack.LevelEnablerFunc(func(lvl logStack.Loglevel) bool { return lvl >= logStack.DebugLevel }),
    },
}
logger := logStack.NewMultiLoger(logs)
```

## Default Logger

The package also provides a default logger that logs to os.Stderr with an InfoLevel. You can retrieve the default logger instance using the DefaultLogger function:

```go
logger := logStack.DefaultLogger()
```

to reset the logger default to use another logging destination use :

```go
logger := logStack.NewLogger(os.Stdout, logStack.InfoLevel)
logStack.ResetDefault(logger)
```

And more
