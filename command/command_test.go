package command

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// At SessionM we use goconvey:

// Finish writing tests so we end up with full test coverage

func TestTransmission(t *testing.T) {
	Convey("Given a new Player", t, func() {
		p := NewPlayer()
		Convey("We can shift down", func() {
			p.ShiftDown()
			So(p.State, ShouldEqual, ShiftingDown)
			err := p.Undo()
			So(err, ShouldBeNil)
			So(p.State, ShouldEqual, Idle)
		})
	})
}
