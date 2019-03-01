intersort
---------

`intersort` sorts slices of arbitrary objects according to the rules of Go's
`fmtsort` package. `fmtsort` is an internal package and thus cannot be imported
directly; however, its behavior is exposed via the `fmt` package when printing
maps. So we can sort arbitrary objects by sticking them in a map, printing it,
and parsing the result.

In other words, this package is an abomination and should not be used for
anything. May the Go maintainers have mercy on my soul.

## Examples:

```go
ints := []int{3, 1, 2}
intersort.Sort(ints) // [1, 2, 3]

strs := []string{"b", "c", "a"}
intersort.Sort(strs) // [a, b, c]

type Point struct {
    X, Y int
}
points := []Point{{2, 1}, {1, 1}, {1, 0}}
intersort.Sort(points) // [{1, 0}, {1, 1}, {2, 1}]
```

You can even sort *differing* types!

```go
objs := []interface{}{3, true, 1, "wat", http.DefaultClient, false}
sort.Sort(intersort.Slice(objs)) // [false, true, 1, 3, wat, &{<nil> <nil> <nil> 0s}]
```

However, the results of this may vary, and in general are unpredictable; see
https://github.com/golang/go/issues/30398.


## Advance praise for intersort:

> I tend to think that exposing the comparison function would be an attractive nuisance. - Ian Lance Taylor

> It was a deliberate decision to keep this implementation private. - Rob Pike

> Sorting [arbitrary objects] is not even possible in the general case. - cznic

> This should not be done. - Rob Pike
