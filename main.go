/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"lolkek/rename/cmd"

	"github.com/eiannone/keyboard"
)

func main() {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	//defer keyboard.Close()
	cmd.Execute()
}
