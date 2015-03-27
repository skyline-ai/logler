# logler

- logs to both stdout/stderror and to [loggly](http://www.loggly.com)
- sample rate for loggly (write only x% of requests to save money), except for emergency log level
- accepts json messages only for loggly, log api for trace
- output: LOG-LEVEL: yyyy/MM/dd hh:mm:ss message
- Trace/Info/Warn/Error/Emergency API

## usage

```go
type JSON map[string]interface{}
func main() {
	logger = logler.New(&logler.Options{
		LogglyToken: "311-234324-2323-234324-2323423",
		Component:   "my-cool-server",
		LogglySampleRate: 25 // (1-100), write only 25% of logs to loggly
	})

	logger.Warn(JSON{"category": "test", "action": "foo", "label": "bar"})
}

```
results in stdout:
```shell
WARNING: 2015/03/27 16:39:30 {"action":"foo","category":"test","label":"bar"}
```

## license
MIT (see LICENSE file)