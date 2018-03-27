package output

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"bytes"
)

func TestNewTable(t *testing.T) {
	Convey("NewTable()", t, func() {
		Convey("Returns a Table", func() {
			So(NewTable(), ShouldNotBeNil)
			So(NewTable(), ShouldHaveSameTypeAs, &Table{})
		})
	})
}

func TestTable(t *testing.T) {
	Convey("Table", t, func() {
		Convey("Is a RecordWriter", func() {
			So(func() {
				var i interface{} = &Table{}
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
				subject := NewTable()
				subject.WriteRecords(buf, records, nil)

				out := buf.String()
				So(out, ShouldContainSubstring, "r1v1")
				So(out, ShouldContainSubstring, "r2v1")
				So(out, ShouldContainSubstring, "r3v1")
			})

			Convey("Writes only from fields", func() {
				buf := &bytes.Buffer{}
				subject := NewTable()
				subject.WriteRecords(buf, records, []string{"key3"})

				out := buf.String()
				So(out, ShouldContainSubstring, "r2v3")
				So(out, ShouldNotContainSubstring, "r2v1")
				So(out, ShouldNotContainSubstring, "r1")
				So(out, ShouldNotContainSubstring, "r3")
			})
		})
	})
}
