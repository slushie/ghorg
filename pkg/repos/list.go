package repos

import (
	"github.com/google/go-github/github"
)

type CompareFunc func(a, b *github.Repository) bool

var DefaultCompareFunc CompareFunc = CompareById

// This type represents a sortable collection of repositories
type List struct {
	Repos   []*github.Repository
	Compare CompareFunc
}

func NewList() *List {
	return &List{
		Repos: make([]*github.Repository, 0),
		Compare: DefaultCompareFunc,
	}
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

func CompareById(a, b *github.Repository) bool {
	return a.GetID() < b.GetID()
}