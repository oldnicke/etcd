package adt_test

import (
	"fmt"

	"oldnicke/etcd/pkg/adt"
)

func Example() {
	ivt := adt.NewIntervalTree()
	ivt.Insert(adt.NewInt64Interval(1, 3), 123)
	ivt.Insert(adt.NewInt64Interval(9, 13), 456)
	ivt.Insert(adt.NewInt64Interval(7, 20), 789)

	rs := ivt.Stab(adt.NewInt64Point(10))
	for _, v := range rs {
		fmt.Printf("Overlapping range: %+v\n", v)
	}
	// output:
	// Overlapping range: &{Ivl:{Begin:7 End:20} Val:789}
	// Overlapping range: &{Ivl:{Begin:9 End:13} Val:456}
}
