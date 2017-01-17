# malwiya
Provide language translation operations using Microsoft translator API with Go.

## Usage
```go
package main

import (
    "log"
    "fmt"
    "github.com/teitei-tk/malwiya"
)

func main() {
    accessToken := "your AccessToken"
    secret := "your secret"
    m := malwiya.New(accessToken, serect)

    from := "en"
    to := "ja"
    text := "I love Golang♡"
    result, err := m.Translate(text, from, to)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(result)
}
```

## LICENSE
MIT
