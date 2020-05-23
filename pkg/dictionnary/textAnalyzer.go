package dictionnary

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type TextAnalyzer interface {
	FeedDictionnary(dictionnary *Dictionnary, reader io.Reader) error
}

type SequentialTextAnalyser struct {
}

func realReadline(reader *bufio.Reader) (*string, error) {
	res := ""
	isPrefix := true
	for isPrefix {
		l, isp, err := reader.ReadLine()
		isPrefix = isp
		if err != nil {
			return nil, err
		}
		res += string(l)
	}
	return &res, nil
}

func (ta *SequentialTextAnalyser) FeedDictionnary(dictionnary *Dictionnary, reader io.Reader) {
	lineReader := bufio.NewReader(reader)
	for {
		curline, err := realReadline(lineReader)
		if err != nil {
			break
		}

		words, err := prepareWords(*curline)

		dictionnary.Feed(words)
	}
}

func prepareWords(line string) ([]string, error) {
	realLine, err := strip(line)
	if err != nil {
		return nil, err
	}
	return strings.Split(realLine, " "), nil
}

func strip(word string) (string, error) {
	var result strings.Builder
	for i := 0; i < len(word); i++ {
		b := word[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}

	if result.Len() == 0 {
		return "", errors.New("no word in line")
	}
	return strings.ToLower(result.String()), nil

}
