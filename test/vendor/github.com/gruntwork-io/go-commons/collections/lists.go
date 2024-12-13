package collections

// ListContainsElement returns true if the given list contains the given element
func ListContainsElement[S ~[]E, E comparable](list S, element any) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}

	return false
}

// RemoveElementFromList returns a copy of the given list with all instances of the given element removed
func RemoveElementFromList[S ~[]E, E comparable](list S, element any) S {
	out := S{}
	for _, item := range list {
		if item != element {
			out = append(out, item)
		}
	}
	return out
}

// MakeCopyOfList will return a new list that is a copy of the given list.
func MakeCopyOfList[S ~[]E, E comparable](list S) S {
	copyOfList := make(S, len(list))
	copy(copyOfList, list)
	return copyOfList
}

// BatchListIntoGroupsOf will group the provided slice into groups of size n, with the last of being truncated to
// the remaining count of elements.  Returns nil if n is <= 0
func BatchListIntoGroupsOf[S ~[]E, E any](slice S, batchSize int) []S {
	if batchSize <= 0 {
		return nil
	}

	// Taken from SliceTricks: https://github.com/golang/go/wiki/SliceTricks#batching-with-minimal-allocation
	// Intuition: We repeatedly slice off batchSize elements from slice and append it to the output, until there
	// is not enough.
	output := []S{}
	for batchSize < len(slice) {
		slice, output = slice[batchSize:], append(output, slice[0:batchSize:batchSize])
	}
	if len(slice) > 0 {
		output = append(output, slice)
	}
	return output
}
