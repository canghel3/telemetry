### Telemetry package

A simple logging package.
Version 1.1.0 is safe for concurrent use.
Version 1.0.1 is usable, but is not concurrency safe.

```bash
go get github.com/Ginger955/telemetry@v1.1.0
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
````

A log can also contain metadata.

```go
//NOTE: re-using an output driver with metadata will result in all messages generated with this driver to contain the given metadata
log.Stdout().WithMetadata(map[any]any{"something":"clean"})
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
log.OutputDriver(ex).Warn().Log("warning")
```
### Levels

<b>Built-in levels:</b>
1. Info
2. Error
3. Warn
4. Debug

```go
log.Stdout().Info().Log("hello world")
```

<b>Extendable</b><br>
Supports addition of self defined levels.

```go
//In-line
log.Stdout().Level(level.Custom("MAJOR")).Log("major level")

//Reusing the level instance
criticalLevel := level.Custom("CRITICAL")
log.Stdout().Level(criticalLevel).Log("critical level")

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
log.Stdout().Settings(filename).Info().Log("with settings overwritten")
```

### Transactions

A transaction can be used to group related logs together.

```go
logTx := log.BeginTx()

logTx.Append(log.Stdout().Info().Msg("first transaction entry"))
logTx.Append(log.File(filename).Error().Msgf("encountered error: %s", "sample error"))
logTx.Log()
```

Transactions also support metadata.

```go
logTx := log.BeginTxWithMetadata(map[any]any{"something":"clean"})

logTx.Append(log.Stdout().Info().Msg("first transaction entry"))
logTx.Append(log.File(filename).Error().Msg("second line is an error"))
logTx.Log()
```



