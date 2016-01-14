package transmission

import (
	"errors"
)

const (
	Idle         = "Idle"
	ShiftingUp   = "Shifting Up"
	ShiftingDown = "Shifting Down"
)

// Represents a transmission which has a current state (only three possible states in the first version) and contains
// a list of all previous states
type Transmission struct {
	State          string
	PreviousStates []string
}

// Returns a new transmission with a call to allocate a new string slice and sets the default transmission state to Idle
func NewTransmission() *Transmission {
	return &Transmission{State: Idle, PreviousStates: make([]string, 0)}
}

// Shifts the transmission up and adds whatever the current state is to the list of previous states
func (t *Transmission) ShiftUp() {
	t.pushCurrentState()
	t.State = ShiftingUp
}

// Shifts the transmission down and adds whatever the current state is to the list of previous states
func (t *Transmission) ShiftDown() {
	t.pushCurrentState()
	t.State = ShiftingDown
}

// Pushes the current state to the list of prior states
func (t *Transmission) pushCurrentState() {
	t.PreviousStates = PushStringToSlice(t.State, t.PreviousStates)
}

// Undoes the last state change the transmission underwent
func (t *Transmission) Undo() error {
	s, err := PopStringFromSlice(t.PreviousStates)
	if err != nil {
		return err
	}
	t.State = s
	return nil
}

// Represents a shift up command. Has a pointer to a transmission
type ShiftUpCommand struct {
	*Transmission
}

// Returns a new shift up command and sets the transmission pointer
func NewShiftUpCommand(t *Transmission) *ShiftUpCommand {
	return &ShiftUpCommand{t}
}

// Represents a shift down command. Has a pointer to a transmission
type ShiftDownCommand struct {
	*Transmission
}

// Returns a new shiftDownCommand and sets the transmission pointer
func NewShiftDownCommand(t *Transmission) *ShiftDownCommand {
	return &ShiftDownCommand{t}
}

// Implements the Execute() part of the Command interface
func (s *ShiftUpCommand) Execute() {
	s.ShiftUp()
}

// Implements the Execute() part of the Command interface
func (s *ShiftDownCommand) Execute() {
	s.ShiftDown()
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
