package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFindSubstrings(t *testing.T) {

	Convey("Given strings 'certain' and 'I'm not certain'...", t, func() {

		s1 := "certain"
		s2 := "I'm not certain"

		Convey("Passing in an empty out array should deliver -1", func() {
			So(findSubstrings(s2, s1, make([]int, 0)), ShouldEqual, -1)
		})

		Convey("Otherwise, should return a new thing containing 1 element", func() {
			out := make([]int, 1)
			So(findSubstrings(s2, s1, out), ShouldEqual, 1)
			So(out[0], ShouldEqual, 8)
		})

	})

}
