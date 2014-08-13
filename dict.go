package main

import (
	"bufio"
	"database/sql"
	"encoding/gob"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"strings"
)

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

	fo, err := os.Create("tokens.gob")
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	words := readPossibleTokens()
	r := bufio.NewWriter(fo)
	w := gob.NewEncoder(r)
	w.Encode(words)
}
