package graphql

import (
	"testing"

	"gitlab.com/Dophin2009/nao/pkg/data"
)

// TestSliceTitles tests the function sliceTitles.
func TestSliceTitles(t *testing.T) {
	point := func(a int) *int {
		return &a
	}

	a := data.Title{
		String: "A", Language: "A", Priority: data.TitlePriorityPrimary}
	b := data.Title{
		String: "B", Language: "B", Priority: data.TitlePrioritySecondary}
	c := data.Title{
		String: "C", Language: "C", Priority: data.TitlePriorityPrimary}
	d := data.Title{
		String: "D", Language: "D", Priority: data.TitlePriorityOther}
	e := data.Title{
		String: "E", Language: "E", Priority: data.TitlePrioritySecondary}
	fullset := []data.Title{a, b, c, d, e}

	cases := []struct {
		name   string
		titles []data.Title
		first  *int
		skip   *int
		res    []*data.Title
	}{
		{"nil:nil:nil", nil, nil, nil, []*data.Title{}},
		{"nil:5:0", nil, point(5), point(0), []*data.Title{}},
		{"5:nil:nil", fullset, nil, nil, []*data.Title{&a, &b, &c, &d, &e}},
		{"5:3:0", fullset, point(3), point(0), []*data.Title{&a, &b, &c}},
		{"5:1:1", fullset, point(1), point(1), []*data.Title{&b}},
		{"5:3:3", fullset, point(3), point(3), []*data.Title{&d, &e}},
		{"5:4:nil", fullset, point(4), nil, []*data.Title{&a, &b, &c, &d}},
		{"5:nil:2", fullset, nil, point(2), []*data.Title{&c, &d, &e}},
		{"5:-1:2", fullset, point(-1), point(2), []*data.Title{&c, &d, &e}},
		{"5:nil:-1", fullset, nil, point(-1), []*data.Title{&a, &b, &c, &d, &e}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := sliceTitles(tc.titles, tc.first, tc.skip)
			if (tc.res == nil) != (res == nil) {
				t.Fatalf("expected nil, but got %v", res)
			}

			if len(tc.res) != len(res) {
				t.Fatalf(
					"expected slice of size %d, but got %d", len(tc.res), len(res))
			}

			for i := range tc.res {
				if *tc.res[i] != *res[i] {
					t.Fatalf("expected %v, but got %v", *tc.res[i], *res[i])
				}
			}
		})
	}
}

// TestCalculatePaginationBounds tests the function calculatePaginationBounds.
func TestCalculatePaginationBounds(t *testing.T) {
	point := func(a int) *int {
		return &a
	}

	cases := []struct {
		name  string
		first *int
		skip  *int
		size  int
		start int
		end   int
	}{
		{"nil:nil:0", nil, nil, 0, 0, 0},
		{"nil:nil:5", nil, nil, 5, 0, 5},
		{"nil:nil:-1", nil, nil, -1, 0, 0},
		{"0:5:5", point(0), point(5), 5, 5, 5},
		{"2:3:5", point(2), point(3), 5, 3, 5},
		{"3:0:5", point(3), point(0), 5, 0, 3},
		{"3:-1:5", point(3), point(-1), 5, 0, 3},
		{"-1:3:5", point(-1), point(3), 5, 3, 5},
		{"3:2:-1", point(3), point(2), -1, 0, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			start, end := calculatePaginationBounds(tc.first, tc.skip, tc.size)
			if start != tc.start {
				t.Fatalf("expected start=%d, but got %d", tc.start, start)
			}
			if end != tc.end {
				t.Fatalf("expected end=%d, but got %d", tc.end, end)
			}
		})
	}
}
