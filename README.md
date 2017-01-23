# malwiya
Provide language translation operations using Microsoft translator Text API with Go.

# Installation
```
go get github.com/teitei-tk/malwiya
```

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

## TODO
* Test code
* Documents

## LICENSE
Apache License, Version 2.0
