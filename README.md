# golang
code base for helper in go language


## How to Use

```go
import "github.com/D074/golang/capture"

capture.InitCapture(&capture.Config{
		OutputPanic: capture.Slack,
		ShowTrace:   true,
		SlackClient: capture.SlackClient{
			WebHookUrl:  "your_web_hook_url",
			UserName:    "Your Sender Name",
			Channel:     "your own channel", 
			MentionUser: []string{}, //optional, ex: "KVGLBJASLN"
			MentionHere: true, //optional,
			Environment: "Your Environment", //optional, ex: prod, beta, stage
		},
	})
```

## Use Go Mux
```go
router := mux.NewRouter()
router.Use(capture.CapturePanic)
```

## Use Http
```go
http.Handle("/foo", capture.CapturePanic(fooHandler))
```

