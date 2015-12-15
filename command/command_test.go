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
		Convey("We can't undo when nothing has been done", func() {
			err := p.Undo()
			So(err, ShouldNotBeNil)
			So(p.State, ShouldEqual, Idle)
		})
		Convey("We can shift down", func() {
			p.ShiftDown()
			So(p.State, ShouldEqual, ShiftingDown)
			err := p.Undo()
			So(err, ShouldBeNil)
			So(p.State, ShouldEqual, Idle)
		})
		Convey("We can shift up", func() {
			p.ShiftUp()
			So(p.State, ShouldEqual, ShiftingUp)
			err := p.Undo()
			So(err, ShouldBeNil)
			So(p.State, ShouldEqual, Idle)
		})
	})
}
