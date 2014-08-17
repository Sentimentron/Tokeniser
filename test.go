package main

import (
	"bufio"
	"database/sql"
	"encoding/gob"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"strings"
    "regexp"
    "bytes"
)

func generateTokenSequenceFromInt(in string, split uint64) []string {
	ret := make([]string, 0)
	last := uint8(0)
	for i := uint8(0); i < 32; i++ {
		if (split & (1 << i)) != 0 {
			ret = append(ret, in[last:i])
			last = i
		}
	}
	ret = append(ret, in[last:])
	return ret
}

func scoreTokenSequence(with map[string]float64, seq []string) (float64, int) {
	ret := 0.0
	for i, s := range seq {
		if v, ok := with[s]; !ok {
			return math.Inf(-1), i
		} else {
			ret += v
		}
	}
	return ret, -1
}

func readPossibleTokens() map[string]float64 {
	fi, err := os.Open("tokens.gob")
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	r := bufio.NewReader(fi)
	d := gob.NewDecoder(r)
	ret := make(map[string]float64)
	d.Decode(&ret)
	return ret
}

func generatePossibleSplitPoints(from map[string]float64, in string, out chan<- uint64) {
    var subStrings []int
    var maxP2 uint8
    var i uint64
    var mask uint64

    // Find split points from dictionary
    for r := range from {
        var c int
        if len(r) <= 2 {
            continue
        }
        for {
            c = findSubstrings(in, r, subStrings)
            if c == -1 {
                subStrings = make([]int, 1 << uint64(len(subStrings)))
            } else {
                break
            }
        }
        for _, i := range subStrings[:c] {
            mask |= (1 << uint8(i))
        }
    }

    // Find maximum possible number
    maxP2 = popCount(mask)
    for i = 0; i < (2 << maxP2); i++ {
        out <- permuteInt(i, mask)
    }
    close(out)
}

func splitToProbableSequence(in string, words map[string]float64) ([]string, float64) {

	if len(in) > 64 {
		return []string{in}, 0.0
	}

	mask := uint64(0)

    tokens := make(chan uint64, mask)
    go generatePossibleSplitPoints(words, in, tokens)

	// Base case: it's already a word
	max, _ := scoreTokenSequence(words, []string{in})
	maxSeq := []string{in}
	inc := uint64(0)

	for i := range tokens {
		seq := generateTokenSequenceFromInt(in, i)
		score, pos := scoreTokenSequence(words, seq)
		if pos != -1 {
			mask |= 1 << uint64(pos)
			inc = countTrailingBits(mask)
		}
		if score > max {
			maxSeq = seq
			max = score
		}
		i &= ^mask
		i += (uint64(1) << inc)
	}
	return maxSeq, max
}

func join(what []string, with string) string {
    var buf bytes.Buffer
    for _, w := range what {
        buf.WriteString(w)
        buf.WriteString(with)
    }
    return buf.String()
}

func main() {

	var cpuprofile = flag.String("cpuprofile", "", "write CPU profile to file")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Compile replacement regexp
    replaceRegex := regexp.MustCompile("[^a-z0-9]")

	wordProbs := readPossibleTokens()

	// Open database
	db, err := sql.Open("sqlite3", "../emotionannotate/spam.sqlite")
	if err != nil {
		panic(err)
	}

	// Read from the input table
	sql := "SELECT document FROM input"
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}

	// Output structure
	tagmap := make(map[string]string)

	defer rows.Close()
	for rows.Next() {
		var doc string
		rows.Scan(&doc)
		words := strings.Split(doc, " ")
		for _, w := range words {
			if len(w) == 0 {
				continue
			}
			if w[0] != '#' {
				continue
			}
			_, ok := tagmap[w]
			if ok {
                continue
            }
			fmt.Print(w)
			fmt.Print(" => ")
            // Need to make the string lower-case and strip punctuation
            w = strings.ToLower(w)
            w = replaceRegex.ReplaceAllString(w, "")
			maxSeq, score := splitToProbableSequence(w, wordProbs)
			fmt.Println(maxSeq, score)
            tagmap[w] = join(maxSeq, " ")
		}
	}

	fo, err := os.Create("hashtags.gob")
    if err != nil {
        panic(err)
    }
    defer fo.Close()
    r := bufio.NewWriter(fo)
    w := gob.NewEncoder(r)
    w.Encode(tagmap)

}
