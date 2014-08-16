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

func generatePossibleSplitPoints(from map[string]float64, in string, out chan uint64) {
    var mask uint64
    var subStrings []int
    var maxP2 uint8
    var i uint64

    // Find split points from dictionary
    for r := range from {
        var c int
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
}

func splitToProbableSequence(in string, words map[string]float64) ([]string, float64) {

	if len(in) > 64 {
		return []string{in}, 0.0
	}

	// Base case: it's already a word
	max, _ := scoreTokenSequence(words, []string{in})
	maxSeq := []string{in}
	mask := uint64(0)
	inc := uint64(0)

	for i := uint64(1); i < uint64(1<<uint64(len(in))); {
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
			fmt.Print(w)
			fmt.Print(" => ")
			maxSeq, score := splitToProbableSequence(w[1:], wordProbs)
			fmt.Println(maxSeq, score)
		}
	}

}
