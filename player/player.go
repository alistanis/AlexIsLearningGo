package player

import (
	"AlexIsLearningGo/car/controllers/shifter"
)

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
	playerCommands *CommandsList
	Shifter        *shiftController
}

// Returns a new player object with a default transmission and commandsList. p.transmission, p.shiftUpCommand, and p.shiftDownCommand all
// contain a pointer to the same transmission (which is why this should be abstracted)
func NewPlayer() *Player {
	commandsList := NewCommandsList()
	shiftController := newShiftController()
	return &Player{playerCommands: commandsList, Shifter: shiftController}
}

// Makes a call to the shift down command via the shiftController
func (p *Player) ShiftDown() {
	p.Shifter.ShiftDownAction(p.playerCommands)
}

// Makes a call to the shift up command via the shiftController
func (p *Player) ShiftUp() {
	p.Shifter.ShiftUpAction(p.playerCommands)
}

// Undoes whatever the last command was and returns an error if there was one
func (p *Player) Undo() error {
	return p.playerCommands.UndoLastCommand()
}

func (p *Player) Redo() error {
	return errors.New("redo not implemented yet")
}
