package collections

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

const (
	DefaultKeyValueStringSliceFormat = "%s=%s"
)

// MergeMaps merges all the maps into one
func MergeMaps[K comparable, V any](mapsToMerge ...map[K]V) map[K]V {
	out := map[K]V{}

	for _, currMap := range mapsToMerge {
		maps.Copy(out, currMap)
	}

	return out
}

// Keys returns the keys for the given map, sorted
func Keys[K constraints.Ordered, V any](m map[K]V) []K {
	out := []K{}

	for key := range m {
		out = append(out, key)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})

	return out
}

// KeyValueStringSlice returns a string slice with key=value items, sorted alphabetically
func KeyValueStringSlice(m map[string]string) []string {
	return KeyValueStringSliceWithFormat(m, DefaultKeyValueStringSliceFormat)
}

// KeyValueStringSliceWithFormat returns a string slice using the specified format, sorted alphabetically.
// The format should consist of at least two format specifiers.
func KeyValueStringSliceWithFormat[K comparable, V any](m map[K]V, format string) []string {
	out := []string{}

	for key, value := range m {
		out = append(out, fmt.Sprintf(format, key, value))
	}

	sort.Strings(out)

	return out
}

// KeyValueStringSliceAsMap converts a string slice with key=value items into a map of slice values. The slice will
// contain more than one item if a key is repeated in the string slice list.
func KeyValueStringSliceAsMap(kvPairs []string) map[string][]string {
	out := make(map[string][]string)
	for _, kvPair := range kvPairs {
		x := strings.Split(kvPair, "=")
		key := x[0]

		var value string
		if len(x) > 1 {
			value = strings.Join(x[1:], "=")
		}

		if _, hasKey := out[key]; hasKey {
			out[key] = append(out[key], value)
		} else {
			out[key] = []string{value}
		}
	}
	return out
}

// MapJoin converts the map to a string type by concatenating the key with the value using the given `mapSep` string, and `sliceSep` string between the slice values.
// For example: `Slice(map[int]string{1: "one", 2: "two"}, "-", ", ")` returns `"1-one, 2-two"`
func MapJoin[M ~map[K]V, K comparable, V any](m M, sliceSep, mapSep string) string {
	list := MapToSlice(m, mapSep)

	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})

	return strings.Join(list, sliceSep)
}

// MapToSlice converts the map to a string slice by concatenating the key with the value using the given `sep` string.
// For example: `Slice(map[int]string{1: "one", 2: "two"}, "-")` returns `[]string{"1-one", "2-two"}`
func MapToSlice[M ~map[K]V, K comparable, V any](m M, sep string) []string {
	var list []string

	for key, val := range m {
		s := fmt.Sprintf("%v%s%v", key, sep, val)
		list = append(list, s)

	}

	return list
}
