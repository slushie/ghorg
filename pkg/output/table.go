package output

import (
	"io"
	"text/tabwriter"
	"fmt"
	"strings"
	"sort"
)

// Table is a RecordWriter that outputs records in tabular format
type Table struct {
	minwidth, tabwidth, padding int
}

var _ RecordWriter = &Table{}

func NewTable() RecordWriter {
	return &Table{
		minwidth: 2,
		tabwidth: 8,
		padding:  2,
	}
}

// Write each record as a line. Note that fields defaults to all fields in the first record.
func (t *Table) WriteRecords(w io.Writer, records []Record, fields []string) error {
	if len(records) == 0 {
		fmt.Fprint(w, "No records found!\n")
		return nil
	}

	tw := tabwriter.NewWriter(w, t.minwidth, t.tabwidth, t.padding, ' ', tabwriter.StripEscape)

	for i, r := range records {
		// Show all fields by default
		if fields == nil {
			var j = 0
			fields = make([]string, len(r))
			for f := range r {
				fields[j] = f
				j ++
			}

			sort.Strings(fields)
		}

		// Output a header row
		if i == 0 {
			fmt.Fprintf(tw, "%s\t\n", strings.Join(fields, "\t"))
			for _, s := range fields {
				fmt.Fprintf(tw, "%s\t", strings.Repeat("-", len(s)))
			}
			fmt.Fprint(tw, "\n")
		}

		// Output the record data
		for _, f := range fields {
			v, _ := r[f] // safely fetch from map
			fmt.Fprintf(tw, "%s\t", v)
		}
		fmt.Fprint(tw, "\n")
	}

	return tw.Flush()
}