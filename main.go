package main

import (
	"fmt"
	"github.com/LeeZXin/zall/cmd"
	"os"
)

func main() {
	app := cmd.NewCliApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
