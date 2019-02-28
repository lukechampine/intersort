package intersort

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"testing"
)

func TestSingleType(t *testing.T) {
	ptrArray := [2]int{2, 1}
	chanArray := [2]chan int{make(chan int), make(chan int)}

	tests := []struct {
		desc   string
		orig   interface{}
		sorted interface{}
	}{
		{
			// "When applicable, nil compares low"
			desc:   "Nil",
			orig:   []*int{&ptrArray[0], (*int)(nil)},
			sorted: []*int{nil, &ptrArray[0]},
		},
		{
			// "ints, floats, and strings order by <"
			desc:   "Ints",
			orig:   []int{3, 2, 1},
			sorted: []int{1, 2, 3},
		},
		{
			// "ints, floats, and strings order by <"
			desc:   "Floats",
			orig:   []float64{3.1, 2.1, 1.1},
			sorted: []float64{1.1, 2.1, 3.1},
		},
		{
			// "ints, floats, and strings order by <"
			desc:   "Strings",
			orig:   []string{"c", "b", "a"},
			sorted: []string{"a", "b", "c"},
		},
		{
			// "NaN compares less than non-NaN floats"
			desc:   "NaN",
			orig:   []float64{3.1, 2.1, math.NaN()},
			sorted: []float64{math.NaN(), 2.1, 3.1},
		},
		{
			// "bool compares false before true"
			desc:   "Bool",
			orig:   []bool{true, false},
			sorted: []bool{false, true},
		},
		{
			// "Complex compares real, then imaginary"
			desc:   "Complex",
			orig:   []complex128{2 + 1i, 1 + 2i},
			sorted: []complex128{1 + 2i, 2 + 1i},
		},
		{
			// "Pointers compare by machine address"
			desc:   "Pointers",
			orig:   []*int{&ptrArray[1], &ptrArray[0]},
			sorted: []*int{&ptrArray[0], &ptrArray[1]},
		},
		{
			// "Channel values compare by machine address"
			desc:   "Channel",
			orig:   []chan int{chanArray[1], chanArray[0]},
			sorted: []chan int{chanArray[0], chanArray[1]},
		},
		{
			// "Structs compare each field in turn"
			desc:   "Structs",
			orig:   []struct{ x, y int }{{1, 0}, {0, 1}},
			sorted: []struct{ x, y int }{{0, 1}, {1, 0}},
		},
		{
			// "Arrays compare each element in turn"
			desc:   "Arrays",
			orig:   [][2]int{{1, 0}, {0, 1}},
			sorted: [][2]int{{0, 1}, {1, 0}},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			Sort(test.orig)
			// can't use reflect.DeepEqual, since NaN != NaN
			if fmt.Sprint(test.orig) != fmt.Sprint(test.sorted) {
				t.Fatalf("expected %#v, got %#v", test.sorted, test.orig)
			}
		})
	}
}

// "Interface values compare first by reflect.Type describing the
// concrete type and then by concrete value as described in the
// previous rules."
func TestMultiType(t *testing.T) {
	ptrArray := [2]int{2, 1}
	chanArray := [2]chan int{make(chan int), make(chan int)}

	// Unfortunately, sorting is unpredictable because it compares the machine
	// addresses of the *reflect.rtype pointers. But we can at least assert that
	// the result contains all sorted subgroups.
	groups := [][]interface{}{
		{(*int)(nil), &ptrArray[0], &ptrArray[1]},
		{1, 2, 3},
		{"a", "b", "c"},
		{false, true},
		{1 + 2i, 2 + 1i},
		{chanArray[0], chanArray[1]},
		{struct{ x, y int }{0, 1}, struct{ x, y int }{1, 0}},
		{[2]int{0, 1}, [2]int{1, 0}},
	}

	var elems Slice
	for _, g := range groups {
		elems = append(elems, g...)
	}
	rand.Shuffle(len(elems), elems.Swap) // nice

	sort.Sort(elems)
	str := fmt.Sprint(elems)
	for _, g := range groups {
		exp := strings.TrimSpace(fmt.Sprintln(g...))
		if !strings.Contains(str, exp) {
			t.Errorf("sorted map should contain %q", exp)
		}
	}
}
