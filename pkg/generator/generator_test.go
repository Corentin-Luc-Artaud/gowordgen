package generator

import (
	"os"
	"runtime/pprof"
	"strconv"
	"testing"

	"github.com/Corentin-Luc-Artaud/gowordgen/pkg/dictionnary"
	"github.com/Corentin-Luc-Artaud/gowordgen/pkg/fetcher"
)

func TestGenerate(t *testing.T) {
	reader, err := fetcher.GeneralFetch("https://www.gutenberg.org/files/1661/1661-0.txt")
	if err != nil {
		t.Error(err)
	}
	dico := dictionnary.NewDictionnary()
	analyser := dictionnary.SequentialTextAnalyser{}
	analyser.FeedDictionnary(dico, reader)

	gen := ParalalGenerator{
		ConcurentsRoutines: 1,
	}

	gen.Use(dico)

	res := gen.Generate(100)
	l := res.Export()
	if len(l) == 0 {
		t.Error("nothing generated")
	}
}

type bench struct {
	gen Generator
	n   int
}

func (ben *bench) oneRun(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if len(ben.gen.Generate(ben.n)) == 0 {
			b.Error("nothing generated")
		}
	}
}

func BenchmarkGenerate(b *testing.B) {
	f, err := os.Create("profiling")
	if err != nil {
		f, err = os.Open("profiling")
		if err != nil {
			b.Error(err)
		}
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	var statMap dictionnary.StatMap

	{
		reader, err := fetcher.GeneralFetch("https://www.gutenberg.org/files/1661/1661-0.txt")
		if err != nil {
			b.Error(err)
		}
		dico := dictionnary.NewDictionnary()
		analyser := dictionnary.SequentialTextAnalyser{}
		analyser.FeedDictionnary(dico, reader)
		statMap = dico.Transitions.Getstatistics()
	}

	for p := 1; p <= 8; p++ {
		gen := &ParalalGenerator{
			ConcurentsRoutines: p,
			statMap:            statMap,
		}
		for n := 1; n < 1000000; n *= 10 {
			ben := bench{
				gen: gen,
				n:   n,
			}

			b.Run("gen"+strconv.Itoa(n)+"_"+strconv.Itoa(p), ben.oneRun)
		}
	}
}
