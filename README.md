# logler

- logs to both stdout/stderror and to loggly.com
- accepts json messages only (map)
- output: LOG-LEVEL: yyyy/MM/dd hh:mm:ss message
- Info/Warn/Error API

## usage

```go
func main() {
	l := logler.New(&logler.Options{
		LogglyToken: "epr3test3pe-1234-test-23e-3test3e",
	})

	l.Warn(&logler.Message{"category": "test", "action": "foo", "label": "bar"})
}

```
results in stdout:
```shell
WARNING: 2015/03/27 16:39:30 {"action":"foo","category":"test","label":"bar"}
```

## license
MIT (see LICENSE file)