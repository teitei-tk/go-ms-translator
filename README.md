# malwiya
Provide language translation operations using Microsoft translator API with Go.

## Usage
```go
package main

import (
    "os"
    "log"
    "fmt"
    "github.com/teitei-tk/malwiya"
)

func main() {
    subscriptionKey := os.Getenv("SUBSCRIPTION_KEY")
    m := malwiya.New(subscriptionKey)

    fromTextLang := "en"
    toTextLang := "ja"
    text := "I love gopher♡"
    result, err := m.Translate(text, fromTextLang, toTextLang)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(result) // 私は gopher♡ が大好き
}
```

## LICENSE
MIT
