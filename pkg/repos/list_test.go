package repos

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"sort"
	"github.com/google/go-github/github"
	"github.com/slushie/ghorg/pkg/output"
	"fmt"
)

func TestListType(t *testing.T) {
	Convey("List type", t, func() {
		Convey("Is sortable", func() {
			Convey("Satisfies sort.Interface", func() {
				var castFunc = func() {
					var subject interface{} = &List{}
					_ = subject.(sort.Interface)
				}

				So(castFunc, ShouldNotPanic)
			})

			Convey(".Add(...repos)", func() {
				Convey("Appends elements to .Repos", func() {
					subject := NewList()
					item := &github.Repository{}

					So(subject.Repos, ShouldHaveLength, 0)
					subject.Add(item)
					So(subject.Repos, ShouldHaveLength, 1)
					So(subject.Repos[0], ShouldEqual, item)
				})
			})

			Convey(".Records(fields) ([]Record, error)", func() {
				var (
					aID int64 = 1
					bID int64 = 2
				)

				subject := NewList()
				subject.Add(
					&github.Repository{ID: &aID},
					&github.Repository{ID: &bID},
				)

				Convey("Calls .Marshal for each repo", func() {
					var calls = 0
					subject.Marshal = func(repo *github.Repository, fields []string) (output.Record, error) {
						calls = calls + 1
						return output.Record{}, nil
					}

					subject.Records(nil)
					So(calls, ShouldEqual, subject.Len())
				})

				Convey("Collects all Records from .Marshal calls", func() {
					records, err := subject.Records(nil)
					So(records, ShouldHaveLength, subject.Len())
					So(err, ShouldBeNil)
				})

				Convey("Stops calling .Marshal on the first error", func() {
					subject.Marshal = func(repo *github.Repository, fields []string) (output.Record, error) {
						return nil, fmt.Errorf("id %d", repo.GetID())
					}

					records, err := subject.Records(nil)
					So(records, ShouldBeNil)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "id 1")
				})

				Convey("Passes fields to .Marshal", func() {
					called := false
					recordFields := []string{"field1"}
					subject.Marshal = func(repo *github.Repository, fields []string) (output.Record, error) {
						called = true
						So(fields, ShouldHaveLength, len(recordFields))
						for i := range recordFields {
							So(fields[i], ShouldEqual, recordFields[i])
						}

						return output.Record{}, nil
					}

					subject.Records(recordFields)
					So(called, ShouldBeTrue)
				})
			})

			Convey(".Len()", func() {
				var tests = []struct {
					in  *List
					out int
				}{
					{&List{Repos: []*github.Repository{}}, 0},
					{&List{Repos: []*github.Repository{nil, nil, nil}}, 3},
				}

				Convey("Returns the length of Repos", func() {
					for _, tt := range tests {
						So(tt.in.Len(), ShouldEqual, tt.out)
					}
				})
			})

			Convey(".Swap(i, j int)", func() {
				var (
					aID int64 = 1
					a         = &github.Repository{ID: &aID}

					bID int64 = 2
					b         = &github.Repository{ID: &bID}
				)

				var tests = []struct {
					in   *List
					i, j int
				}{
					{&List{Repos: []*github.Repository{a, b}}, 0, 1},
					{&List{Repos: []*github.Repository{a, b}}, 1, 1},
				}

				Convey("Swaps items at i and j", func() {
					for _, tt := range tests {
						var (
							i = tt.in.Repos[tt.i]
							j = tt.in.Repos[tt.j]
						)

						tt.in.Swap(tt.i, tt.j)
						So(i, ShouldEqual, tt.in.Repos[tt.j])
						So(j, ShouldEqual, tt.in.Repos[tt.i])
					}
				})

			})

			Convey(".Less(i, j int)", func() {
				var (
					aID int64 = 1
					bID int64 = 2
				)
				var (
					a = &github.Repository{ID: &aID}
					b = &github.Repository{ID: &bID}
				)

				Convey("Returns true if the item at i should sort before j", func() {
					var tests = []struct {
						c    CompareFunc
						r    []*github.Repository
						i, j int
						v    bool
					}{
						{DefaultCompareFunc, []*github.Repository{a, b}, 0, 1, true},
						{DefaultCompareFunc, []*github.Repository{a, b}, 1, 0, false},
						{DefaultCompareFunc, []*github.Repository{a, a}, 0, 1, false},
					}

					for _, tt := range tests {
						var subject = &List{Repos: tt.r, Compare: tt.c}
						So(subject.Less(tt.i, tt.j), ShouldEqual, tt.v)
					}
				})

				Convey("Calls the Compare func", func() {
					var called = false
					var subject = &List{
						Repos: []*github.Repository{nil, nil},
						Compare: func(a, b *github.Repository) bool {
							called = true
							return true
						},
					}

					v := subject.Less(0, 1)
					So(called, ShouldEqual, true)
					So(v, ShouldEqual, true)
				})
			})
		})
	})
}

func TestNewList(t *testing.T) {
	Convey("NewList() *List", t, func() {

		Convey("Returns a non-nil *List", func() {
			So(NewList(), ShouldNotBeNil)
		})

		Convey("Makes a slice for Repos", func() {
			So(NewList().Repos, ShouldNotBeNil)
		})

		Convey("Assigns a CompareFunc", func() {
			So(NewList().Repos, ShouldNotBeNil)
		})
	})
}

func TestCompareByID(t *testing.T) {
	Convey("CompareByID(a, b *Repository) bool", t, func() {
		var (
			aID int64 = 1
			bID int64 = 2
		)
		var (
			a = &github.Repository{ID: &aID}
			b = &github.Repository{ID: &bID}
		)

		Convey("Compares with GetID()", func() {
			var tests = []struct {
				a, b *github.Repository
				v    bool
			}{
				{a, b, true},
				{b, a, false},
				{a, nil, false}, // GetID returns 0 for nil repos
				{a, a, false},
			}

			for _, tt := range tests {
				So(CompareByID(tt.a, tt.b), ShouldEqual, tt.v)
			}
		})
	})
}
