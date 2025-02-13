package commands

type CommandHandler func(args ...string) error
type Commands map[string]CommandHandler

var commands Commands

func Exist(name string) bool {
	return commands[name] != nil
}

func Get(name string) CommandHandler {
	return commands[name]
}

func RegisterCommand(commandName string, handler CommandHandler) {
	if commands == nil {
		commands = make(Commands)
	}

	commands[commandName] = handler
}

func RegisterCommands() {
	RegisterCommand("ls", ls)
	RegisterCommand("cls", cls)
	RegisterCommand("mk", mk)
	RegisterCommand("mkdir", mkdir)
	RegisterCommand("rm", rm)
	RegisterCommand("task", task)
	RegisterCommand("cd", cd)
	RegisterCommand("mget", mget)
}
