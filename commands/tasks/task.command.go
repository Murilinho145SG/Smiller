package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"smiller/commands/ic"
	"smiller/lines"
)

type handler func(args ...string) error
type handlerRegistry map[string]*tasks

type tasks struct {
	Handler     handler
	Description string
}

var TasksRegistries handlerRegistry

func init() {
	if TasksRegistries == nil {
		TasksRegistries = make(handlerRegistry)
	}

	register("-l", "lists all created automations", list)
	register("test", "test commands", test)
	register("start", "start tasks saved in task.json", start)
}

func register(n, d string, h handler) {
	t := &tasks{
		Handler:     h,
		Description: d,
	}

	TasksRegistries[n] = t
}

func list(args ...string) error {
	b, err := os.ReadFile("./task.json")
	if err != nil {
		return err
	}

	var v map[string]interface{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if len(v) == 0 {
		return errors.New("tasks is empty")
	}

	for k := range v {
		fmt.Println(k)
	}

	return nil
}

type Task struct {
	Name string
	*Params
}

type Params struct {
	Commands []string
	Files    []string
}

func GetTasks() ([]*Task, error) {
	b, err := os.ReadFile("./task.json")
	if err != nil {
		return nil, err
	}

	var v map[string]map[string][]string
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}

	var ts []*Task
	for k, p := range v {
		t := &Task{
			Name:   k,
			Params: &Params{},
		}

		for k2, m := range p {
			if k2 == "command" {
				t.Commands = append(t.Commands, m...)
				continue
			}

			t.Files = append(t.Files, m...)
		}

		ts = append(ts, t)
	}

	return ts, nil
}

func GetTask(taskName string) (*Task, error) {
	b, err := os.ReadFile("./task.json")
	if err != nil {
		return nil, err
	}

	var v map[string]map[string][]string
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}

	tm := v[taskName]
	if tm == nil {
		return nil, errors.New("invalid task")
	}

	t := &Task{
		Name:   taskName,
		Params: &Params{},
	}

	if tm["command"] != nil {
		t.Commands = tm["command"]
	}

	if tm["files"] != nil {
		t.Files = tm["files"]
	}

	return t, nil
}

func test(args ...string) error {
	if args[1] == "test" {

	}

	return nil
}

func readFiles(taskName string) (map[string][]byte, error) {
	de, err := os.ReadDir("./tasks/" + taskName)
	if err != nil {
		return nil, err
	}

	files := make(map[string][]byte)

	for _, entry := range de {
		if entry.IsDir() {
			return nil, errors.New("is not support folders")
		}

		b, err := os.ReadFile("./tasks/" + taskName + "/" + entry.Name())
		if err != nil {
			return nil, err
		}

		files[entry.Name()] = b
	}

	return files, nil
}

func start(args ...string) error {
	t, err := GetTask(args[2])
	if err != nil {
		return err
	}

	if len(t.Commands) > 0 {
		for _, v := range t.Commands {
			p := lines.Parser(v)
			c := exec.Command(p[0], p[1:]...)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Dir = ic.ActualDir

			if err := c.Run(); err != nil {
				return err
			}
		}
	}

	if len(t.Files) > 0 {
		files, err := readFiles(t.Name)
		if err != nil {
			return errors.New("not have files")
		}

		for _, v := range t.Files {
			p := filepath.Join(ic.ActualDir, v)
			fmt.Println("creating", v)
			f, err := os.Create(p)
			if err != nil {
				return err
			}
			defer f.Close()

			file := files[v]
			if file == nil {
				return errors.New("files list not correspond")
			}

			f.Write(file)
		}
	}

	return nil
}
