package graphql

import (
	"testing"

	"github.com/Dophin2009/nao/pkg/data/models"
)

// TestSliceTitles tests the function sliceTitles.
func TestSliceTitles(t *testing.T) {
	point := func(a int) *int {
		return &a
	}

	a := models.Title{
		String: "A", Language: "A", Priority: models.TitlePriorityPrimary}
	b := models.Title{
		String: "B", Language: "B", Priority: models.TitlePrioritySecondary}
	c := models.Title{
		String: "C", Language: "C", Priority: models.TitlePriorityPrimary}
	d := models.Title{
		String: "D", Language: "D", Priority: models.TitlePriorityOther}
	e := models.Title{
		String: "E", Language: "E", Priority: models.TitlePrioritySecondary}
	fullset := []models.Title{a, b, c, d, e}

	cases := []struct {
		name   string
		titles []models.Title
		first  *int
		skip   *int
		res    []*models.Title
	}{
		{"nil:nil:nil", nil, nil, nil, []*models.Title{}},
		{"nil:5:0", nil, point(5), point(0), []*models.Title{}},
		{"5:nil:nil", fullset, nil, nil, []*models.Title{&a, &b, &c, &d, &e}},
		{"5:3:0", fullset, point(3), point(0), []*models.Title{&a, &b, &c}},
		{"5:1:1", fullset, point(1), point(1), []*models.Title{&b}},
		{"5:3:3", fullset, point(3), point(3), []*models.Title{&d, &e}},
		{"5:4:nil", fullset, point(4), nil, []*models.Title{&a, &b, &c, &d}},
		{"5:nil:2", fullset, nil, point(2), []*models.Title{&c, &d, &e}},
		{"5:-1:2", fullset, point(-1), point(2), []*models.Title{&c, &d, &e}},
		{"5:nil:-1", fullset, nil, point(-1), []*models.Title{&a, &b, &c, &d, &e}},
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
