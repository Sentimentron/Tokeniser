package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math"
	"os"
	"runtime/pprof"
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

	var maxSeq []string
	var mask uint64
	words := readPossibleTokens()
	max := math.Inf(-1)

	for i := uint64(0); i < (1 << 11); i++ {
		if i&mask != 0 {
			continue
		}
		seq := generateTokenSequenceFromInt("geniusbaby", i)
		score, pos := scoreTokenSequence(words, seq)
		if pos != -1 {
			mask |= (1 << uint64(pos))
		}
		if score > max {
			maxSeq = seq
			max = score
		}
	}

	fmt.Println(max, maxSeq)
	//possibleLexemes := generatePossibleTokens("geniusbaby")
	//fmt.Println(possibleLexemes)
	//possibleCombinations := generatePossibleTokenCombinations(possibleLexemes)
	//fmt.Println(possibleCombinations)
	//fmt.Println(generateCoveringLexemes("geniusbaby", possibleLexemes))
	// fmt.Println(words)
}
