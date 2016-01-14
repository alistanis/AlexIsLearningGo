package command

import (
	"errors"
	"sync"
)

// this could implement the ability to have a global list of all commands run - think of the possibilities here
var (
	allCommands = NewCommandsList()
)

// Things to consider:
// implement redo
// implement observer/notifier pattern when a command is performed (in terms of games think
// about achievements)

type Command interface {
	Execute()
	Undo() error
}

// Keeps a list of Command interfaces (interfaces are always a pointer type)
type CommandsList struct {
	List  []Command
	mutex sync.Mutex
}

// Returns a new command list with a call to make([]Command, 0) (this is Go's dynamic allocation command, similar to malloc)
func NewCommandsList() *CommandsList {
	return &CommandsList{List: make([]Command, 0), mutex: sync.Mutex{}}
}

// Pushes a command onto the list. Protected by a mutex in case this list of commands is global and can be referenced
// from multiple contexts/goroutines
func (c *CommandsList) PushCommand(command Command) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.List = append(c.List, command)
}

// Undoes the last command and removes it from the list. Returns an error when this is impossible. Mutex protected
func (c *CommandsList) UndoLastCommand() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	command, err := c.PopLastCommand()
	if err != nil {
		return err
	}
	return command.Undo()
}

// Pops the last command off of the list. Not mutex protected, but is threadsafe when called from UndoLastCommand()
func (c *CommandsList) PopLastCommand() (Command, error) {
	if len(c.List) == 0 {
		return nil, errors.New("Nothing to pop off the slice")
	}
	var command Command
	command, c.List = c.List[len(c.List)-1], c.List[:len(c.List)-1]
	return command, nil
}
