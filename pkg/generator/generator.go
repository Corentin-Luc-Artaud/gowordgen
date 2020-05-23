package generator

import (
	"log"
	"math/rand"
	"strings"

	"github.com/Corentin-Luc-Artaud/gowordgen/pkg/dictionnary"
	"github.com/Corentin-Luc-Artaud/gowordgen/pkg/utils"
)

type Generator interface {
	Use(dictionnary *dictionnary.Dictionnary)
	Generate(n int) utils.Set
}

type ParalalGenerator struct {
	ConcurentsRoutines int
	statMap            dictionnary.StatMap
}

func (generator *ParalalGenerator) Use(dictionnary *dictionnary.Dictionnary) {
	generator.statMap = dictionnary.Transitions.Getstatistics()
}

func (generator *ParalalGenerator) generateOne(r *rand.Rand) (string, error) {
	state := dictionnary.Start + dictionnary.Start
	res := strings.Builder{}
	new := dictionnary.Start
	for new != dictionnary.End {
		random := r.Float64() * 100
		new, err := generator.statMap.GetCorrespondingTransition(state, random)
		if err != nil {
			return "", err
		}
		if new == dictionnary.End {
			break
		}
		if _, err := res.WriteString(new); err != nil {
			return "", err
		}
		state = (state + new)[1:]
	}
	return res.String(), nil
}

func (generator *ParalalGenerator) work(n int, resChan chan []string, seed int64) {
	r := rand.New(rand.NewSource(seed))
	res := make([]string, n)
	for i := 0; i < n; i++ {
		new, err := generator.generateOne(r)
		if err != nil {
			log.Fatal(err)
		}
		res[i] = new
	}
	resChan <- res
}

func (generator *ParalalGenerator) Generate(n int) utils.Set {
	coWorker := generator.ConcurentsRoutines
	if n < generator.ConcurentsRoutines {
		coWorker = 1
	}

	load := n / coWorker

	resChan := make(chan []string)

	for i := 0; i < coWorker; i++ {
		go generator.work(load, resChan, rand.Int63())
	}

	res := utils.Set{}
	//get results from all workers
	for i := 0; i < coWorker; i++ {
		res.AddUnexistants(<-resChan)
	}
	return res
}
