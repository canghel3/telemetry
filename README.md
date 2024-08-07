### Telemetry package

A simple logging package. A major refactoring (v1.1.0) is underway. Version 1.0.1 is usable, but is not concurrency safe.

```bash
go get -u github.com/Ginger955/telemetry@v1.0.1
```

### Output
<b>Built-in outputs:</b><br>
1. File logging

```go
//Logging in-line
log.File(filename).Info().Log([]byte("hello, world!"))

//Reusing the same output destination
toFile := log.File(filename)

toFile.Info().Log("foo")
toFile.Error().Logf("encountered error: %s", "sample error")
```

2. Stdout logging

```go
//Logging in-line
log.Stdout().Info().Log("hello, world!")

//Reusing the same output destination
stdout := log.Stdout()

stdout.Info().Log("foo")
stdout.Error().Logf("encountered error: %s", "sample error")
```

A log can also contain metadata.

```go
log.Stdout().Metadata(map[any]any{"something":"clean"})
```

<b>Extendable</b> <br>
Supports addition of custom output drivers for logging to any custom implementation.
```go
type example struct {
	msg []byte
}

func (e example) Write(b []byte) (int, error) {
	e.msg = b
	return len(b), nil
}


ex := example{}
log.OutputDriver(ex).Warn().Log([]byte("warning"))
```
### Levels

<b>Built-in levels:</b>
1. NoLevel
2. Error
3. Warn
4. Info
5. Debug

```go
log.Stdout().Info().Log([]byte("hello world"))
```

<b>Extendable</b><br>
Supports addition of self defined levels.

```go
//In-line
log.Stdout().Level(level.Custom("MAJOR")).Log([]byte("major level"))

//Reusing the level instance
criticalLevel := level.Custom("CRITICAL")
log.Stdout().Level(criticalLevel).Log([]byte("critical level"))

```

### Configuration

The log outputs can be customized using a configuration file. Configuration is limited to timestamp formatting and enabling/disabling implicit output formatting. <br>

```json
{
  "formatting": {
    "log" : {
      "disabled": false,
      "field_order":  {
        "timestamp": 1,
        "level": 2,
        "metadata": 3,
        "buffer": 4
      },
      "timestamp": "2024-03-02"
    },
    "transaction" : {
      "disabled": false,
      "field_order":  {
        "timestamp": 1,
        "level": 2,
        "metadata": 3,
        "buffer": 4
      },
      "timestamp": "2024-03-02"
    }
  }
}
```

Each logger instance can be modified using a different configuration file.

```go
log.Stdout().Settings(filename).Info().Log([]byte("with settings overwritten"))
```

### Transactions

A transaction can be used to group related logs together.

```go
logTx := log.BeginTx()

logTx.Append(log.Stdout().Info().Msg([]byte("first transaction entry")))
logTx.Append(log.File(filename).Error().Msg([]byte("second line is an error")))
logTx.Log()
```

Transactions also support metadata.

```go
logTx := log.BeginTxWithMetadata(map[any]any{"something":"clean"})

logTx.Append(log.Stdout().Info().Msg([]byte("first transaction entry")))
logTx.Append(log.File(filename).Error().Msg([]byte("second line is an error")))
logTx.Log()
```



