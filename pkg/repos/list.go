package repos

import (
	"github.com/google/go-github/github"
	"github.com/slushie/ghorg/pkg/output"
	"fmt"
	"sort"
)

type CompareFunc func(a, b *github.Repository) bool
type MarshalFunc func(repo *github.Repository, fields []string) (output.Record, error)

var DefaultCompareFunc CompareFunc = CompareByID
var DefaultMarshalFunc MarshalFunc = MarshalID

// This type represents a sortable collection of repositories
type List struct {
	Repos   []*github.Repository
	Compare CompareFunc
	Marshal MarshalFunc
}

func NewList() *List {
	return &List{
		// sane default capacity -- ptrs are small.
		Repos:   make([]*github.Repository, 0, 10),
		Compare: DefaultCompareFunc,
		Marshal: DefaultMarshalFunc,
	}
}

// Append repos... to the list
func (l *List) Add(repos ...*github.Repository) {
	l.Repos = append(l.Repos, repos...)
}

func (l *List) Records(fields []string) ([]output.Record, error) {
	records := make([]output.Record, l.Len())
	for i, repo := range l.Repos {
		if rec, err := l.Marshal(repo, fields); err != nil {
			return nil, err
		} else {
			records[i] = rec
		}
	}
	return records, nil
}

// Sort data in-place
func (l *List) Sort() {
	sort.Sort(l)
}

// Sort data in-place in reverse order
func (l *List) SortReverse() {
	sort.Sort(sort.Reverse(l))
}

func (l *List) Len() int {
	return len(l.Repos)
}

func (l *List) Swap(i, j int) {
	l.Repos[i], l.Repos[j] = l.Repos[j], l.Repos[i]
}

func (l *List) Less(i, j int) bool {
	return l.Compare(l.Repos[i], l.Repos[j])
}

func CompareByID(a, b *github.Repository) bool {
	return a.GetID() < b.GetID()
}

func MarshalID(repo *github.Repository, fields []string) (output.Record, error) {
	return output.Record{
		"id": fmt.Sprint(repo.GetID()),
	}, nil
}