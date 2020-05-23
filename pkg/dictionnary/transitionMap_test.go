package dictionnary_test

import (
	"testing"

	"github.com/Corentin-Luc-Artaud/gowordgen/pkg/dictionnary"
	"github.com/Corentin-Luc-Artaud/gowordgen/pkg/utils"
)

func TestMerge(t *testing.T) {
	one := dictionnary.TransitionMap{
		"a": map[string]int64{
			"b": int64(1),
		},
	}
	two := dictionnary.TransitionMap{
		"a": map[string]int64{
			"b": int64(1),
		},
		"b": map[string]int64{
			"c": int64(1),
		},
	}

	expected := dictionnary.TransitionMap{
		"a": map[string]int64{
			"b": int64(2),
		},
		"b": map[string]int64{
			"c": int64(1),
		},
	}

	one.Merge(two)

	for previous, possible := range expected {
		onePossible, presents := one[previous]
		if !presents {
			t.Error(previous, " not presents in the resulted map")
		}
		for k, v := range possible {
			oneV, presents := onePossible[k]
			if !presents {
				t.Error(k, " not presents for ", previous)
			}
			if oneV != v {
				t.Error(k, " has a bad value for ", previous)
			}
		}
	}
}

func TestFeed(t *testing.T) {
	words := []string{"bonjour"}
	dic := dictionnary.Dictionnary{
		Words:       utils.NewSet(),
		Transitions: dictionnary.TransitionMap{},
	}
	dic.Feed(words)

	if !dic.Words.ContainsAll(words) {
		t.Error("word not added")
	}

	if dic.Transitions[dictionnary.Start+dictionnary.Start]["b"] != 1 {
		t.Error("bad value for ||[b]")
	}

	if dic.Transitions[dictionnary.Start+"b"]["o"] != 1 {
		t.Error("bad value for |b[o]")
	}

	if dic.Transitions["bo"]["n"] != 1 {
		t.Error("bad value for bo[n]")
	}

	if dic.Transitions["on"]["j"] != 1 {
		t.Error("bad value for on[j]")
	}

	if dic.Transitions["nj"]["o"] != 1 {
		t.Error("bad value for nj[o]")
	}

	if dic.Transitions["jo"]["u"] != 1 {
		t.Error("bad value for jo[u]")
	}
	if dic.Transitions["ou"]["r"] != 1 {
		t.Error("bad value for ou[r]")
	}

	if dic.Transitions["ur"][dictionnary.End] != 1 {
		t.Error("bad value for ur[/]")
	}
}
