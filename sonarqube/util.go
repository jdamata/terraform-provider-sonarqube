package sonarqube

import (
	"reflect"
	"sort"
)

// Checks if two string slices are equal, optionally ignoring ordering
func stringSlicesEqual(a, b []string, ignoreOrder bool) bool {
	if ignoreOrder {
		sort.Slice(a, func(i, j int) bool {
			return a[i] < a[j]
		})
		sort.Slice(b, func(i, j int) bool {
			return b[i] < b[j]
		})
	}

	return reflect.DeepEqual(a, b)
}
