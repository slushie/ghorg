package output

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"bytes"
)

func TestNewJson(t *testing.T) {
	Convey("NewJson()", t, func() {
		Convey("Returns a Json", func() {
			So(NewJson(), ShouldNotBeNil)
			So(NewJson(), ShouldHaveSameTypeAs, &Json{})
		})
	})
}


func TestJson(t *testing.T) {
	Convey("Json", t, func() {
		Convey("Is a RecordWriter", func() {
			So(func() {
				var i interface{} = &Json{}
				_ = i.(RecordWriter)
			}, ShouldNotPanic)
		})

		Convey(".WriteRecords()", func() {
			var records = []Record{
				{"key1": "r1v1", "key2": "r1v2"},
				{"key1": "r2v1", "key2": "r2v2", "key3": "r2v3"},
				{"key1": "r3v1", "key2": "r3v2", "key4": "r3v4"},
			}

			Convey("Writes all records", func() {
				buf := &bytes.Buffer{}
				subject := NewJson()
				subject.WriteRecords(buf, records, nil)

				out := buf.String()
				So(out, ShouldContainSubstring, "r1v1")
				So(out, ShouldContainSubstring, "r2v1")
				So(out, ShouldContainSubstring, "r3v1")
			})

			// Json doesn't support fields
		})
	})
}