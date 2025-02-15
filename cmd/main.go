package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"smiller/commands"
	"smiller/lines"
	"strings"
)

type Config struct {
	User string `json:"user"`
}

func ReadConfig() (*Config, error) {
	b, err := os.ReadFile("./smiller.config.json")
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	c, err := ReadConfig()
	if err != nil {
		fmt.Print("Error for read config file")
		os.Exit(1)
	}

	commands.RegisterCommands()
	dir, err := os.Getwd()
	if err != nil {
		fmt.Print("\033[31mError:", err.Error())
		return
	}

	commands.ActualDir = dir

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\033[34m"+c.User, " \033[35m|Smiller>\033[0m ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Exiting...")
			break
		}

		args := lines.Parser(input)
		if len(args) == 0 {
			continue
		}

		if commands.Exist(args[0]) {
			ch := commands.Get(args[0])
			if err := ch(args...); err != nil {
				fmt.Println("\033[31mError:\033[0m", err.Error())
			}
			continue
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Dir = commands.ActualDir

		if err := cmd.Run(); err != nil {
			fmt.Println("Error for exec command: executable file not found")
		}
	}
}
