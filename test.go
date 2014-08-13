package main

import (
	"database/sql"
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

func generatePossibleTokenCombinations(in []string) [][]string {

	ret := make([][]string, 0)
	for i := 1; i <= len(in); i++ {
		// i is a possible substring length
		for j := 0; j <= len(in)-i; j++ {
			// j is a possible starting position
			substring := in[j : i+j]
			ret = append(ret, substring)
		}

	}

	return ret

}

func generatePossibleTokens(in string) []string {

	ret := make([]string, 0)

	for i := 1; i <= len(in); i++ {
		// i is a possible substring length
		for j := 0; j <= len(in)-i; j++ {
			// j is a possible starting position
			substring := in[j : i+j]
			ret = append(ret, substring)
		}

	}

	return ret

}

func readPossibleTokens() map[string]float64 {

	ret := make(map[string]float64)
	total := 0
	tmp := make(map[string]int)

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
			// Ignore hashtags, at-mentions
			if w[0] == '#' || w[0] == '@' {
				continue
			}
			tmp[w]++
			total++
		}
	}

	for i := range tmp {
		ret[i] = math.Log(float64(tmp[i]) / float64(total))
	}

	return ret

}

func tokensCover(substrings []string, str string) bool {

	candidate := strings.Join(substrings, "")
	return candidate == str

}

/*func generateCoveringLexemes(of string, from []string) [][]string {

	var str bytes.Buffer // Contains the complete reference string
	var buf bytes.Buffer
	var ins int // Append pointer

	// List of things that we can choose
	canChoose := make([]bool, len(from))
	chosen := make([]bool, len(from)) // List of things we have chosen

	// Copy original string to buffer
	str.WriteString(of)

	// Create output
	ret := make([][]string, 0)

	// Termination condition
	numChosen := func(s []bool) int {
		ret := 0
		for _, c := range s {
			if c {
				ret++
			}
		}
		return ret
	}

	for i := range canChoose {
		canChoose[i] = true
	}

	for {
		buf.Reset()
		// Initialise the list of things chosen
		for i := range chosen {
			chosen[i] = false
			canChoose[i] = true
		}
		// Reset pointer to the head of the string
		ins = 0

		for i, c := range canChoose {
			if !c {
				continue
			}
			if chosen[i] {
				continue
			}
			f := from[i]
			l := len(f)

			fmt.Println(ins, l, buf.String())
			if ins+l > len(of) {
				break
			}
			fmt.Println(of[ins : ins+l])
			if of[ins:ins+l] == f {
				// This is a covering candidate
				chosen[i] = true
				ins += l
				buf.WriteString(f)
			}
		}

		n := numChosen(chosen)
		if n == 0 {
			break
		} else {
			retBuf := make([]string, n)
			retCnt := 0
			for i, c := range chosen {
				if !c {
					continue
				}
				retBuf[retCnt] = from[i]
				canChoose[i] = false
				retCnt++
			}
			ret = append(ret, retBuf)
		}
	}

	return ret
}*/

/*func generateCoveringLexemes(of string, from []string) [][]string {

	var comp bytes.Buffer
	var buf bytes.Buffer
	var tmp bytes.Buffer
	canChoose := make([]bool, len(from))
	chose := make([]bool, len(from))
	comp.WriteString(of)

	ret := make([][]string, 0)

	for i := range canChoose {
		canChoose[i] = true
	}

	// Start with the empty string
	for {
		// Reset the list of things we chose
		for i := range chose {
			chose[i] = false
		}
		// Choose the next string which maintains
		// covering
		for i, n := range canChoose {
			if !n {
				continue
			}
			// Copy buf to tmp
			tmp.Reset()
			tmp.Write(buf.Bytes())
			// Append candidate string
			tmp.WriteString(from[i])
			// Check if the buffer still covers the string
			tmpBytes := tmp.Bytes()
			cmpBytes := comp.Bytes()
			matches := true
			for i, b := range tmpBytes {
				if b != cmpBytes[i] {
					matches = false
					break
				}
			}
			if matches {
				// This can be a token, so we can't choose it any more
				canChoose[i] = false
				// Record that we can add this to the finished set
				chose[i] = true
				buf.WriteString(from[i])
			}
		}
		// Check to see if the entire buffer matches
		matches := true
		cmpBytes := comp.Bytes()
		if buf.Len() != len(cmpBytes) {
			matches = false
		} else {
			for i, b := range buf.Bytes() {
				if b != cmpBytes[i] {
					matches = false
					break
				}
			}
		}
		if matches {
			// Emit this string
			retBuf := make([]string, 0)
			for i, c := range chose {
				if !c {
					continue
				}
				retBuf = append(retBuf, from[i])
			}
			ret = append(ret, retBuf)
		}
		// Termination condition: nothing got chosen
		somethingChosen := false
		for _, c := range chose {
			if c {
				somethingChosen = true
				break
			}
		}
		if !somethingChosen {
			break
		}
	}
	return ret

}
*/
func generatePossibleTokenizations(forString string, fromSubset []string) [][]string {

	ret := make([][]string, 0)
	for i, _ := range fromSubset {
		// i is a possible starting position
		for j := 0; j < len(fromSubset)-i; j++ {
			// j is a possible length
			substrings := fromSubset[j : i+j+1]
			fmt.Println(substrings)
			if tokensCover(substrings, forString) {
				ret = append(ret, substrings)
			}
		}
	}

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

	var maxSeq uint64
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
			maxSeq = i
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
