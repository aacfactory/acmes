package main

import (
	"fmt"
	"github.com/aacfactory/acmes/internal/command"
)

func main() {
	err := command.Run()
	if err != nil {
		fmt.Println(err)
	}
}
