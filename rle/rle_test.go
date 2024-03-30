package rle_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/makyo/gogol/rle"
)

func TestMarshal(t *testing.T) {
	Convey("Given a starting field", t, func() {
		f := &rle.RLEField{
			Width:  3,
			Height: 3,
			Field: [][]bool{
				[]bool{false, true, false},
				[]bool{false, false, true},
				[]bool{true, true, true},
			},
			Name:    "Test",
			Origin:  "Tester",
			Survive: []int{2, 3},
			Born:    []int{3},
		}

		Convey("When it is marshalled", func() {
			result := f.Marshal()

			Convey("The result should be an RLE file's contents", func() {
				So(result, ShouldEqual, `#N Test
#O Tester
#R 0  0
x = 3, y = 3, rule = B3/S23
bo$2bo$3o!
`)
			})
		})
	})
}

func TestUnmarshal(t *testing.T) {
	Convey("When unmarshalling the contents of an RLE file", t, func() {
		f, err := rle.Unmarshal(`#N Sample
#O Tester
#C Comment
#CXRLE Pos=0,-1377 Gen=3480106827776
x = 12, y = 12, rule = B3/S23
4b2o$4b2o2$4b4o$3bobo2bob2o$3bo2bobob2o$2ob2o3bo$2obo4bo$4b4o2$6b2o$6b
2o!`)
		Convey("It should not error", func() {
			So(err, ShouldBeNil)
		})

		Convey("It sets metadata properly", func() {
			So(f.Name, ShouldEqual, "Sample")
			So(f.Origin, ShouldEqual, "Tester")
			So(f.Comments, ShouldEqual, []string{"Comment"})
			So(f.ExtendedRLEData, ShouldEqual, []string{"Pos=0,-1377 Gen=3480106827776"})
			So(f.Top, ShouldEqual, 0)
			So(f.Left, ShouldEqual, 0)
		})

		Convey("It parses the headers", func() {
			So(f.Width, ShouldEqual, 12)
			So(f.Height, ShouldEqual, 12)
			So(f.Survive, ShouldEqual, []int{2, 3})
			So(f.Born, ShouldEqual, []int{3})
		})

		Convey("It parses the rule block", func() {
			So(f.Field, ShouldEqual, [][]bool{
				[]bool{false, false, false, false, true, true, false, false, false, false, false, false},
				[]bool{false, false, false, false, true, true, false, false, false, false, false, false},
				[]bool{false, false, false, false, false, false, false, false, false, false, false, false},
				[]bool{false, false, false, false, true, true, true, true, false, false, false, false},
				[]bool{false, false, false, true, false, true, false, false, true, false, true, true},
				[]bool{false, false, false, true, false, false, true, false, true, false, true, true},
				[]bool{true, true, false, true, true, false, false, false, true, false, false, false},
				[]bool{true, true, false, true, false, false, false, false, true, false, false, false},
				[]bool{false, false, false, false, true, true, true, true, false, false, false, false},
				[]bool{false, false, false, false, false, false, false, false, false, false, false, false},
				[]bool{false, false, false, false, false, false, true, true, false, false, false, false},
				[]bool{false, false, false, false, false, false, true, true, false, false, false, false},
			})
		})
	})
}
