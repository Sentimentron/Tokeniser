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

func splitToProbableSequence(in string, words map[string]float64) ([]string, float64) {

	var maxSeq []string
	var mask uint64
	max := math.Inf(-1)
	for i := uint64(0); i < (1 << uint64(len(in))); i++ {
		seq := generateTokenSequenceFromInt(in, i)
		score, pos := scoreTokenSequence(words, seq)
		//	fmt.Println(seq, score, pos)
		if pos != -1 {
			mask |= (1 << uint64(pos))
		}
		if score > max {
			maxSeq = seq
			max = score
		}
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
