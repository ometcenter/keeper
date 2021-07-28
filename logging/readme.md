# logging

применение `sentry`:
```go

package main

import (
	"fmt"

	"github.com/ometcenter/keeper/config"
	"github.com/ometcenter/keeper/logging"
)

func main() {

	c := &config.LoggerConfig{
		Level: int(logging.DebugLevel),
	}

	s, err := logging.NewSentryLog("http://6b200a96a9c54653b5395ba417abbf3a@localhost:9000/3")
	if err != nil {
		fmt.Printf("Ошибка при создании экземпляра sentry: %v\n", err)
	}
	logging.InitLog(c, s)

	logging.Impl.Info("message")

}
```

результат:
```go
C:\projects\keeper-stack\keeper>go run main.go
[Sentry] 2021/07/28 18:50:04 Integration installed: ContextifyFrames
[Sentry] 2021/07/28 18:50:04 Integration installed: Environment
[Sentry] 2021/07/28 18:50:04 Integration installed: Modules
[Sentry] 2021/07/28 18:50:04 Integration installed: IgnoreErrors
http://6b200a96a9c54653b5395ba417abbf3a@localhost:9000/3
[INF] 2021/07/28 - 18:50:04 | message
```