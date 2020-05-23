package dictionnary

import "github.com/Corentin-Luc-Artaud/gowordgen/pkg/utils"

type Dictionnary struct {
	Words       utils.Set     `json:"words"`
	Transitions TransitionMap `json:"transitions"`
}

func NewDictionnary() *Dictionnary {
	return &Dictionnary{
		Words:       utils.Set{},
		Transitions: TransitionMap{},
	}
}

type MarshallableDictionnary struct {
	Words       []string      `json:"words"`
	Transitions TransitionMap `json:"transitions"`
}

func (dic *Dictionnary) GetTransitionMap() TransitionMap {
	return dic.Transitions
}

func (dic *Dictionnary) Feed(words []string) {

	for _, word := range dic.Words.AddUnexistants(words) {
		dic.Transitions.Feed(word)
	}
}

func (dic *Dictionnary) Marshallable() *MarshallableDictionnary {
	return &MarshallableDictionnary{
		Transitions: dic.Transitions,
		Words:       dic.Words.Export(),
	}
}

func (mdic *MarshallableDictionnary) ToReal() *Dictionnary {
	return &Dictionnary{
		Transitions: mdic.Transitions,
		Words:       utils.CreateSet(mdic.Words),
	}
}
