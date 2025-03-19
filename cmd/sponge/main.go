// Package main sponge is a basic development framework that integrates code auto generation,
// Gin and GRPC, a microservice framework. it is easy to build a complete project from development
// to deployment, just fill in the business logic code on the generated template code, greatly improved
// development efficiency and reduced development difficulty, the use of Go can also be "low-code development".
package main

import (
	"fmt"
	"github.com/go-dev-frame/sponge/cmd/sponge/global"
	"os"
	"path/filepath"

	"github.com/go-dev-frame/sponge/cmd/sponge/commands"
	"github.com/go-dev-frame/sponge/cmd/sponge/commands/generate"
)

func main() {
	err := generate.Init()
	if err != nil {
		fmt.Printf("\n    %v\n\n", err)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	global.Root = filepath.Dir(filepath.Dir(dir))
	global.Path = dir

	rootCMD := commands.NewRootCMD()
	if err = rootCMD.Execute(); err != nil {
		rootCMD.PrintErrln("Error:", err)
		os.Exit(1)
	}
}
