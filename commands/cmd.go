package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"smiller/commands/ic"
	"smiller/commands/tasks"
	"smiller/utils"
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

	path := ActualDir + "\\" + value
	f, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !f.IsDir() {
		return errors.New("is not a dir")
	}

	ActualDir += "\\" + value
	return nil
}

func mget(args ...string) error {
	if len(args) == 1 {
		utils.System("use mget help to list then commands")
		return nil
	}

	value := args[1]
	if value == "" {
		return errors.New("invalid arguments")
	}

	if value == "help" {
		help(map[string]string{
			"[url]": "make a get requisition for the url",
		})
		return nil
	}

	if !strings.HasPrefix(value, "http") {
		return errors.New("invalid url")
	}

	r, err := http.Get(value)
	if err != nil {
		return err
	}

	fmt.Println(r.Status, r.Request.Proto)
	for k, v := range r.Header {
		fmt.Println(k+":", v)
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	fmt.Println("\n" + string(b))

	return nil
}

func mpost(args ...string) error {
	if len(args) == 1 {
		utils.System("use mpost help to list then commands")
		return nil
	}

	value := args[1]
	if value == "" {
		return errors.New("invalid arguments")
	}

	if value == "help" {
		help(map[string]string{
			"[url]": "make a post requisition for the url",
		})
		return nil
	}

	if !strings.HasPrefix(value, "http") {
		return errors.New("invalid url")
	}

	body := args[2]
	if body == "" {
		return errors.New("no body present")
	}

	r, err := http.Post(value, "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		return err
	}

	fmt.Println(r.Status, r.Request.Proto)
	for k, v := range r.Header {
		fmt.Println(k+":", v)
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	fmt.Println("\n" + string(b))

	return nil
}

func help(values map[string]string) {
	for k, v := range values {
		fmt.Println("	\033[33m"+k, "\033[0m"+v)
	}
}

func task(args ...string) error {
	if len(args) == 1 {
		utils.System("use task help to list the commands")
		return nil
	}

	var v string
	if len(args) > 1 {
		v = args[1]
	}

	if v == "help" && len(args) == 2 {
		handlerDesc := make(map[string]string)
		for k, t := range tasks.TasksRegistries {
			handlerDesc[k] = t.Description
		}

		help(handlerDesc)
		return nil
	}

	ic.ActualDir = ActualDir

	t := tasks.TasksRegistries[v]

	if t == nil {
		return errors.New("invalid params, use help for the commands")
	}

	err := t.Handler(args...)
	if err != nil {
		return err
	}

	return nil
}
