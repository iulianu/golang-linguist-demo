package main

import (
	"fmt"
	"bytes"
	"strings"
	"regexp"
	"errors"
	"strconv"
	"io"
)

type Phrase struct {
	tokenCount uint32
	tokens string
}

type PhrasePair struct {
	source Phrase
	target Phrase
	memo string
}

func (this Phrase) ToMungedForm() string {
	return fmt.Sprintf("%d %s", this.tokenCount, this.tokens)
}

func (this PhrasePair) ToMungedForm() string {
	return this.source.ToMungedForm() +
		" $ " +
		this.target.ToMungedForm() +
		" & " +
		this.memo +
		" |";
}

func fromUnmungedString(text []byte) (PhrasePair, error) {
	//fmt.Printf("Unmunged: %s\n", text)
	re := regexp.MustCompile("PP (\\d+) (\\d+) (.*) \\$ (.*) & (.*)")
	subm := re.FindSubmatch(text)
	if subm == nil {
		return PhrasePair{}, errors.New("Phrase pair syntax error")
	}
	pair := PhrasePair{}
	pair.source = Phrase{}
	pair.target = Phrase{}
	var err error
	var cnt int
	cnt, err = strconv.Atoi(string(subm[1]))
	if err != nil { return PhrasePair{}, err }
	pair.source.tokenCount = uint32(cnt)
	cnt, err = strconv.Atoi(string(subm[2]))
	if err != nil { return PhrasePair{}, err }
	pair.target.tokenCount = uint32(cnt)
	pair.source.tokens = string(subm[3])
	pair.target.tokens = string(subm[4])
	pair.memo = string(subm[5])
	return pair, nil
}

func ParsePhrasePairStream(text []byte) ([]PhrasePair, error) {
	re := regexp.MustCompile("(?:\\A|\\s+)PP \\d+ \\d+")
	pairs := make([]PhrasePair, 0)
	prevTokenStart := 0
	for prevTokenStart <= len(text) {
		remainingSlice := text[prevTokenStart:]
		locs := re.FindAllIndex(remainingSlice, 2)
		var end int
		if locs == nil || len(locs) < 2 {
			end = len(remainingSlice)
		} else {
			end = locs[1][0] + 1 // TODO not sane
		}
		token := bytes.TrimSpace(remainingSlice[0:end])
		if len(token) > 0 {
			pp, err := fromUnmungedString(token)
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, pp)
			prevTokenStart += end
		} else {
			prevTokenStart += end + 1
		}
	}
	return pairs, nil
}

func ParsePhrasePairReader(w io.Writer, r io.RuneReader) error {
	re := regexp.MustCompile("(?:\\A|\\s+)PP \\d+ \\d+")
	//prevTokenStart := 0
	re.FindReaderIndex(r)
	return errors.New("Phrase pair syntax error")
}

func MungedPairs(pairs []PhrasePair) string {
	textPairs := make([]string, 0)
	for _, v := range pairs {
		textPairs = append(textPairs, v.ToMungedForm())
	}
	return strings.Join(textPairs, " ")
}
