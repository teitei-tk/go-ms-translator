package main

import (
	"fmt"
	"os"

	"github.com/teitei-tk/malwiya"
)

func main() {
	m := malwiya.New(os.Getenv("SUBSCRIPTION_KEY"))
	fmt.Println(m.Translate("i love gopher♡", "en", "ja"))
}
