package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"smiller/lines"
	"strings"
)

var ActualDir string

func ls(args ...string) error {
	dir := ActualDir
	if len(args) == 2 {
		dir = args[1]
	}

	de, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var values []any
	for _, entry := range de {

		if entry.Type().IsDir() {
			values = append(values, fmt.Sprintf("\033[36m%s\033[0m", entry.Name()))
			continue
		}

		values = append(values, entry.Name())
	}

	if len(values) == 0 {
		return errors.New("this path not have files")
	}

	fmt.Println(values...)
	return nil
}

func cls(args ...string) error {
	if len(args) > 1 {
		return errors.New("invalid input")
	}

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()

	return nil
}

func mk(args ...string) error {
	fP := filepath.Join(ActualDir, args[1])

	f, err := os.Create(fP)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

func mkdir(args ...string) error {
	fP := filepath.Join(ActualDir, args[1])
	err := os.MkdirAll(fP, 0644)
	if err != nil {
		return err
	}

	return nil
}

func rm(args ...string) error {
	errInvalidArguments := errors.New("invalid arguments!")
	if len(args) > 3 {
		return errInvalidArguments
	}

	param := args[1]
	path := filepath.Join(ActualDir, args[2])

	if param == "-a" {
		return os.RemoveAll(path)
	}

	if param == "-o" {
		return os.Remove(path)
	}

	return errInvalidArguments
}

func cd(args ...string) error {
	if len(args) < 2 {
		return errors.New("invalid arguments")
	}

	value := args[1]
	if value == "" {
		return errors.New("invalid params")
	}

	if value == ".." {
		i := strings.LastIndex(ActualDir, "\\")
		ActualDir = ActualDir[:i]
		fmt.Println(ActualDir)
		return nil
	}

	f, err := os.Stat(value)
	if err != nil {
		return err
	}

	if !f.IsDir() {
		return errors.New("is not a dir")
	}

	s := path.Join(value)
	fmt.Println(s)
	ActualDir += "\\" + s
	fmt.Println(ActualDir)

	return nil
}

func task(args ...string) error {
	if len(args) > 3 {
		return errors.New("invalid params")
	}

	var tasks map[string]interface{}
	b, err := os.ReadFile("./task.json")
	if err != nil {
		return err
	}

	taskName := args[1]

	err = json.Unmarshal(b, &tasks)
	if err != nil {
		return err
	}

	taskInfo := tasks[taskName]
	if taskInfo == nil {
		return errors.New("non-existent task")
	}

	taskMap, ok := taskInfo.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid structure for task '%s'", taskName)
	}

	commands, ok := taskMap["command"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid structure for 'command' in task '%s'", taskName)
	}

	var cmdStrings []string
	for _, cmd := range commands {
		if str, ok := cmd.(string); ok {
			cmdStrings = append(cmdStrings, str)
		} else {
			return fmt.Errorf("invalid type in 'command' array for task '%s'", taskName)
		}
	}

	fmt.Println("Commands:", cmdStrings)

	for _, cmd := range cmdStrings {
		v := args[2]
		if v != "" {
			s := strings.Split(v, ":")
			varName := s[0]
			newName := s[1]
			cmd = strings.ReplaceAll(cmd, varName, newName)
		}

		values := lines.Parser(cmd)
		c := exec.Command(values[0], values[1:]...)
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		c.Dir = ActualDir

		if err := c.Run(); err != nil {
			return err
		}

		fmt.Println("Executing command:", cmd)
	}

	return nil
}
