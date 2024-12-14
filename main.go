package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/go-bumbu/todo-app/app/cmd"
)

func main() {
	cmd.Execute()
}

var _ = spew.Dump
