package types

/***** IntValuesList *********************************************************/

type IntValuesList []int

// Swap exchanges the elements at indexes i and j.
func (l IntValuesList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Swap exchanges the elements at indexes i and j.
func (l IntValuesList) Len() int {
	return len(l)
}

// Swap exchanges the elements at indexes i and j.
func (l IntValuesList) Less(i, j int) bool {
	return l[i] < l[j]
}

/***** SampleSetsList ********************************************************/

type SampleSetsList []*SampleSet

// Swap exchanges the elements at indexes i and j.
func (l SampleSetsList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Swap exchanges the elements at indexes i and j.
func (l SampleSetsList) Len() int {
	return len(l)
}

// Swap exchanges the elements at indexes i and j.
func (l SampleSetsList) Less(i, j int) bool {
	return l[i].Less(l[j])
}

/***** SlicesList ************************************************************/

type SlicesList []*Slice

// Swap exchanges the elements at indexes i and j.
func (l SlicesList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Swap exchanges the elements at indexes i and j.
func (l SlicesList) Len() int {
	return len(l)
}

// Swap exchanges the elements at indexes i and j.
func (l SlicesList) Less(i, j int) bool {
	return l[i].Less(l[j])
}
