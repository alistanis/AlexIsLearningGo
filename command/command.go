package command

import (
	"errors"
	"sync"
)

const (
	Idle         = "Idle"
	ShiftingUp   = "Shifting Up"
	ShiftingDown = "Shifting Down"
)

// this could implement the ability to have a global list of all commands run - think of the possibilities here
var (
	allCommands = NewCommandsList()
)

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

// Things to consider:
// implement redo
// implement observer/notifier pattern when a command is performed (in terms of games think
// about achievements)

type Command interface {
	Execute()
	Undo() error
}

// Pushes a string to a slice and returns the (potentially) newly allocated slice.
// Slices in Go are always a pointer type, and calls to append can change the address of the slice, thus losing the
// reference, so any function that takes a slice and mutates it must also return a slice
func PushStringToSlice(s string, slc []string) []string {
	slc = append(slc, s)
	return slc
}

// Pops the last string from a slice and returns an error if an error occured
func PopStringFromSlice(slc []string) (string, error) {
	// In Go, on a slice or map type, calling len will return either the length of the slice, which could actually be
	// 0, but len(slc) == 0 will also return 0 if slc is nil
	if len(slc) == 0 {
		return "", errors.New("Nothing to pop off the string slice")
	}
	var s string
	s, slc = slc[len(slc)-1], slc[:len(slc)-1]
	return s, nil
}

// Represents a transmission which has a current state (only three possible states in the first version) and contains
// a list of all previous states
type transmission struct {
	State          string
	PreviousStates []string
}

// Returns a new transmission with a call to allocate a new string slice and sets the default transmission state to Idle
func newTransmission() *transmission {
	return &transmission{State: Idle, PreviousStates: make([]string, 0)}
}

// Shifts the transmission up and adds whatever the current state is to the list of previous states
func (t *transmission) ShiftUp() {
	t.pushCurrentState()
	t.State = ShiftingUp
}

// Shifts the transmission down and adds whatever the current state is to the list of previous states
func (t *transmission) ShiftDown() {
	t.pushCurrentState()
	t.State = ShiftingDown
}

// Pushes the current state to the list of prior states
func (t *transmission) pushCurrentState() {
	t.PreviousStates = PushStringToSlice(t.State, t.PreviousStates)
}

// Undoes the last state change the transmission underwent
func (t *transmission) Undo() error {
	s, err := PopStringFromSlice(t.PreviousStates)
	if err != nil {
		return err
	}
	t.State = s
	return nil
}

// Things to consider:
// Should the player have a direct reference to the transmission to get state?
// (we should probably have something like a shifter(controller) that
// manages the interaction between the player and the transmission,
// rather than having the player have access to the shift
// commands manually. Think about MVC.

// Additionally - the Player object in this case would be responsible for more than just a simple state machine -
// - the player would be able to turn and shift at the same time, but those can still be executed as commands/undo

// The Player struct contains pointers to a transmission, shiftUpCommand, shiftDownCommand (those should be moved to a
// shiftController or shifter struct) and contains a list of the previous commands run by the player.
type Player struct {
	*transmission
	*shiftUpCommand
	*shiftDownCommand
	playerCommands *CommandsList
}

// Returns a new player object with a default transmission and commandsList. p.transmission, p.shiftUpCommand, and p.shiftDownCommand all
// contain a pointer to the same transmission (which is why this should be abstracted)
func NewPlayer() *Player {
	t := newTransmission()
	commandsList := NewCommandsList()
	return &Player{transmission: t, playerCommands: commandsList, shiftUpCommand: newShiftUpCommand(t), shiftDownCommand: newShiftDownCommand(t)}
}

// Makes a call to the shift down command
func (p *Player) ShiftDown() {
	p.shiftDownCommand.Execute()
	p.playerCommands.PushCommand(p.shiftDownCommand)
}

// Makes a call to the shift up command
func (p *Player) ShiftUp() {
	p.shiftUpCommand.Execute()
	p.playerCommands.PushCommand(p.shiftUpCommand)
}

// Undoes whatever the last command was and returns an error if there was one
func (p *Player) Undo() error {
	return p.playerCommands.UndoLastCommand()
}

// Gets the current shift state (actually returns p.transmission.State) (this should be abstracted)
func (p *Player) GetShiftState() string {
	return p.State
}

// Represents a shift up command. Has a pointer to a transmission
type shiftUpCommand struct {
	*transmission
}

// Returns a new shift up command and sets the transmission pointer
func newShiftUpCommand(t *transmission) *shiftUpCommand {
	return &shiftUpCommand{t}
}

// Represents a shift down command. Has a pointer to a transmission
type shiftDownCommand struct {
	*transmission
}

// Returns a new shiftDownCommand and sets the transmission pointer
func newShiftDownCommand(t *transmission) *shiftDownCommand {
	return &shiftDownCommand{t}
}

// Implements the Execute() part of the Command interface
func (s *shiftUpCommand) Execute() {
	s.ShiftUp()
}

// Implements the Execute() part of the Command interface
func (s *shiftDownCommand) Execute() {
	s.ShiftDown()
}
