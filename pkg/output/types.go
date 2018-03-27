package output

import "io"

type Record map[string]string

type RecordWriter interface {
	// WriteRecords writes each record to w in order. If fields is non-nil,
	// it should provide a hint to the field ordering.
	WriteRecords(w io.Writer, records []Record, fields []string) error
}
