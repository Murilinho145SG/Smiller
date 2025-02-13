package tasks

type Handler func(args ...string) error

type Tasks struct {
	Name        string
	Handler     Handler
	Description string
}

func Create(args ...string) error {

	return nil
}

func Flags(args ...string) error {

	return nil
}

func Exec(args ...string) error {

	return nil
}
