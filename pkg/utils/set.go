package utils

//Set is an implementation of set for string
type Set map[string]struct{}

//Add add word to the set if not present, return true if correctly added
func (set *Set) Add(word string) bool {
	if _, present := (*set)[word]; present {
		return false
	}
	(*set)[word] = struct{}{}
	return true
}

//Contains return true if the set contains word
func (set *Set) Contains(word string) bool {
	_, present := (*set)[word]
	return present
}

func (set *Set) ContainsAll(words []string) bool {
	for _, word := range words {
		if !set.Contains(word) {
			return false
		}
	}
	return true
}

func (set *Set) AddUnexistants(words []string) []string {
	toAdd := make([]string, 0, len(words))
	for _, word := range words {
		if set.Add(word) {
			toAdd = append(toAdd, word)
		}
	}
	return toAdd
}

//Export exports the presnts words
func (set *Set) Export() []string {
	keys := make([]string, 0, len(*set))
	for k := range *set {
		keys = append(keys, k)
	}
	return keys
}

//CreateSet create a set from a slice of string
func CreateSet(values []string) Set {
	set := Set(make(map[string]struct{}))
	for _, v := range values {
		set.Add(v)
	}
	return set
}

func NewSet() Set {
	set := Set(make(map[string]struct{}))
	return set
}
