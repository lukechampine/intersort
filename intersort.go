package intersort // import "lukechampine.com/intersort"

import (
	"fmt"
	"reflect"
	"sort"
)

func less(x, y interface{}) bool {
	justX := fmt.Sprint(map[interface{}]struct{}{x: struct{}{}})
	justY := fmt.Sprint(map[interface{}]struct{}{y: struct{}{}})
	return fmt.Sprint(map[interface{}]struct{}{
		x: struct{}{},
		y: struct{}{},
	}) == fmt.Sprintf("map[%v %v]", justX[4:len(justX)-1], justY[4:len(justY)-1])
}

// Slice implements sort.Interface for arbitrary objects, according to the map
// ordering of the fmt package.
type Slice []interface{}

func (is Slice) Len() int           { return len(is) }
func (is Slice) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }
func (is Slice) Less(i, j int) bool { return less(is[i], is[j]) }

// Sort sorts arbitrary objects according to the map ordering of the fmt
// package.
func Sort(slice interface{}) {
	val := reflect.ValueOf(slice)
	if val.Type().Kind() != reflect.Slice {
		panic("intersort: cannot sort non-slice type")
	}
	sort.Slice(slice, func(i, j int) bool {
		return less(val.Index(i).Interface(), val.Index(j).Interface())
	})
}
