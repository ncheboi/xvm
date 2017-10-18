package main

import (
	"fmt"

	"."
)

func main() {
	_, err := xvm.NewRun()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("hello")
	}
}
