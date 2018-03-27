package output

import (
	"io"
	"encoding/json"
)

type Json struct {
	prefix, indent string
}

func NewJson() RecordWriter {
	return &Json{indent: "  "}
}

func (j *Json) WriteRecords(w io.Writer, records []Record, fields []string) error {
	e := json.NewEncoder(w)
	e.SetIndent("", j.indent)

	return e.Encode(records)
}
