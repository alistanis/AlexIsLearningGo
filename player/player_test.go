package player

import (
	. "AlexIsLearningGo/car/transmission"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// At SessionM we use goconvey:

// Finish writing tests so we end up with full test coverage

func TestTransmission(t *testing.T) {
	Convey("Given a new Player", t, func() {
		p := NewPlayer()
		Convey("We can't undo when nothing has been done", func() {
			err := p.Undo()
			So(err, ShouldNotBeNil)
			So(p.Shifter.GetTransmissionState(), ShouldEqual, Idle)
		})
		Convey("We can shift down", func() {
			p.ShiftDown()
			So(p.Shifter.GetTransmissionState(), ShouldEqual, ShiftingDown)
			err := p.Undo()
			So(err, ShouldBeNil)
			So(p.Shifter.GetTransmissionState(), ShouldEqual, Idle)
		})
		Convey("We can shift up", func() {
			p.ShiftUp()
			So(p.Shifter.GetTransmissionState(), ShouldEqual, ShiftingUp)
			err := p.Undo()
			So(err, ShouldBeNil)
			So(p.Shifter.GetTransmissionState(), ShouldEqual, Idle)
		})
		Convey("We can redo something that has been undone", func() {
			p.ShiftUp()
			So(p.Shifter.GetTransmissionState(), ShouldEqual, ShiftingUp)
			err := p.Undo()
			So(err, ShouldBeNil)
			So(p.Shifter.GetTransmissionState(), ShouldEqual, Idle)
			err = p.Redo()
			So(err, ShouldNotBeNil)
			So(p.Shifter.GetTransmissionState(), ShouldEqual, Idle)
		})
	})
}
