package shift_controller

import (
	"AlexIsLearningGo/car/transmission"
)

// Manages interactions between Player and Transmission
type shiftController struct {
	*Transmission
	*ShiftDownCommand
	*ShiftUpCommand
}

// Returns a new shiftController with pointers to a transmission, shiftUpCommand and shiftDownCommand.
// Things to consider: once we have multiple players each shifting their own transmission, should those interactions all be handled by a single shiftController, or should each player/transmission pair get it's own shiftController?
func newShiftController() *shiftController {
	t := NewTransmission()
	d := NewShiftDownCommand(t)
	u := NewShiftUpCommand(t)
	return &shiftController{t, d, u}
}

// Handles shifting a transmission up for a player
func (s *shiftController) ShiftUpAction(c *CommandsList) {
	s.shiftUpCommand.Execute()

	// don't know if this is the best approach, but I don't think the alternative - returning the shiftCommand as a "success" value and having the Player's shiftUp method push that onto the slice - is really any better
	c.PushCommand(s.shiftUpCommand)
}

// Handles shifting a transmission down for a player
func (s *shiftController) ShiftDownAction(c *CommandsList) {
	s.shiftDownCommand.Execute()
	c.PushCommand(s.shiftDownCommand)
}

// Handles undoing the last shift action for a player
func (s *shiftController) UndoLastAction() {

}

// Handles redoing the last action undone for a player
func (s *shiftController) RedoLastUndoAction() error {
	return errors.New("redo not implemented yet")
}

// Returns this shiftController's transmission's State
func (s *shiftController) GetTransmissionState() string {
	return s.transmission.State
}
