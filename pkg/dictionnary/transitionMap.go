package dictionnary

import (
	"errors"
)

const (
	Start = "|"
	End   = "/"
)

//StatMap represents the probabilities for each transitions given a state
//a,b => [{c, 1}, {d, 99}]
//the sum of all transition for a given state should be 100
type StatMap map[string]map[string]float64

//TransitionMap represents the occurences for each transition
// a,b => [{c, 999}, {d, 99999}]
type TransitionMap map[string]map[string]int64

//TransitionMap methods

func (tmap *TransitionMap) incrementTransition(state string, transition string) {
	if _, presents := (*tmap)[state]; !presents {
		(*tmap)[state] = map[string]int64{}
	}
	if _, presents := (*tmap)[state][transition]; !presents {
		(*tmap)[state][transition] = 0
	}
	(*tmap)[state][transition]++
}

func (tmap *TransitionMap) Feed(word string) {
	state := Start + Start
	for _, char := range word {
		tmap.incrementTransition(state, string(char))
		state = (state + string(char))[1:]
	}
	tmap.incrementTransition(state, End)
}

//Merge merge a transition map into the current
func (tmap *TransitionMap) Merge(o TransitionMap) {
	for transition, occurs := range o {
		currentOccurs, present := (*tmap)[transition]
		if !present {
			(*tmap)[transition] = occurs
			return
		}

		for k, v := range occurs {
			currentValue, present := currentOccurs[k]
			if !present {
				currentOccurs[k] = v
			}
			currentOccurs[k] = currentValue + v
		}
	}
}

func computeTotal(transitions map[string]int64) int64 {
	var res int64 = 0
	for _, v := range transitions {
		res += v
	}
	return res
}

func computeStats(transitions map[string]int64, total float64) map[string]float64 {
	stats := map[string]float64{}
	for k, v := range transitions {
		proba := float64(v*100) / total
		if proba > 100 {
			//error
		}
		stats[k] = float64(proba)
	}
	return stats
}

//Getstatistics build the corresponding statistique map
func (tmap *TransitionMap) Getstatistics() StatMap {
	res := StatMap{}
	for from, transitions := range *tmap {
		total := computeTotal(transitions)
		res[from] = computeStats(transitions, float64(total))
	}
	//res.computeOrders()
	return res
}

//GetCorrespondingTransition return the transition corresponding to the random number and the transition,
//ThreadSafe (readonly)
func (statMap *StatMap) GetCorrespondingTransition(from string, random float64) (string, error) {
	var cumulativeStat float64 = 0
	for k, v := range (*statMap)[from] {
		cumulativeStat += v
		if random < cumulativeStat {
			return k, nil
		}
	}
	return "", errors.New("not valid probabilities")
}
