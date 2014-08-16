package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPopCount(t *testing.T) {
	Convey("Given something with 8 1 bits...", t, func() {
		num := uint64(0xFF)
		Convey("Result should be 8...", func() {
			So(popCount(num), ShouldEqual, 8)
		})
	})
}

func TestPermuteInt(t *testing.T) {

	Convey("With a 100101 mask...", t, func() {
		mask := uint64(0x29)
		So(permuteInt(0, mask), ShouldEqual, 0)
		So(permuteInt(1, mask), ShouldEqual, 1)
		So(permuteInt(2, mask), ShouldEqual, 8)
		So(permuteInt(3, mask), ShouldEqual, 9)
		So(permuteInt(4, mask), ShouldEqual, 32)
		So(permuteInt(5, mask), ShouldEqual, 33)
		So(permuteInt(6, mask), ShouldEqual, 40)
		So(permuteInt(7, mask), ShouldEqual, 41)
	})

}
