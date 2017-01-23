package main

import (
	"fmt"
	"os"

	"github.com/teitei-tk/malwiya"
)

func main() {
	m := malwiya.New(os.Getenv("SUBSCRIPTION_KEY"))
	r, err := m.TrasnlateArray([]string{"I", "love", "you"}, "en", "ja")
	if err != nil {
		panic(err)
	}

	fmt.Println(r)
}
