package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Corentin-Luc-Artaud/gowordgen/pkg/dictionnary"
	"github.com/Corentin-Luc-Artaud/gowordgen/pkg/fetcher"
	"github.com/Corentin-Luc-Artaud/gowordgen/pkg/generator"
	"github.com/Corentin-Luc-Artaud/gowordgen/pkg/utils"
)

func createDictionnary(dicPath string, books string) {
	dico := dictionnary.Dictionnary{
		Words:       utils.Set{},
		Transitions: dictionnary.TransitionMap{},
	}
	feedDictionnary(books, &dico)
	saveDictionnary(&dico, dicPath)
}

func updateDictionnary(dicPath string, books string) {
	dico := loadDictionnary(dicPath)
	feedDictionnary(books, dico)
	saveDictionnary(dico, dicPath)
}

func feedDictionnary(book string, dico *dictionnary.Dictionnary) {
	bookPaths := strings.Split(book, string(os.PathListSeparator)+string(os.PathListSeparator))
	analyzer := dictionnary.SequentialTextAnalyser{}
	for _, p := range bookPaths {
		reader, err := fetcher.GeneralFetch(p)
		if err != nil {
			continue
		}
		analyzer.FeedDictionnary(dico, reader)
	}
}

func loadDictionnary(path string) *dictionnary.Dictionnary {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	dico := dictionnary.MarshallableDictionnary{
		Transitions: dictionnary.TransitionMap{},
		Words:       []string{},
	}
	json.Unmarshal(data, &dico)
	return dico.ToReal()
}

func saveDictionnary(dico *dictionnary.Dictionnary, path string) {
	data, err := json.MarshalIndent(dico.Marshallable(), "", "	")
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		log.Fatal(err)
	}
}

func computeBook(books string, dicPath string) {
	if books == "" {
		log.Fatal("books should contains at least one path try -h for help")
	}
	if dicPath == "" {
		log.Fatal("please specify a dictionnary path, try -h for help")
	}

	info, err := os.Stat(dicPath)
	if err != nil {
		createDictionnary(dicPath, books)
	} else if info.IsDir() {
		log.Fatal("dictionnary is a directory")
	} else {
		updateDictionnary(dicPath, books)
	}
}

func generateWords(dicPath string, n int, threads int) {
	if dicPath == "" {
		log.Fatal("please specify a dictionnary path, try -h for help")
	}

	dico := loadDictionnary(dicPath)
	gen := generator.ParalalGenerator{
		ConcurentsRoutines: threads,
	}

	gen.Use(dico)

	//startTime := time.Now()
	res := gen.Generate(n)
	//duration := time.Since(startTime)

	for _, w := range res.Export() {
		fmt.Println(w)
	}

	//println()
	//println(duration.Seconds())
}

func main() {
	generate := flag.Int("generate", 0, "number of words, if 0 compute books for the given dictionnary")
	dicPath := flag.String("dico", "", "path of the dictionnary to use")
	books := flag.String("books", "", "book paths to feed the dictionnary (separated by :: ) could be http/https url")
	threads := flag.Int("p", 1, "number of threads used to generates words")

	flag.Parse()

	if *generate > 0 {
		generateWords(*dicPath, *generate, *threads)
	} else {
		computeBook(*books, *dicPath)
	}
}
