package main

import ()

func findSubstrings(srcString, expectedSubstr string, out []int) int {
	outCtr := 0
	for i := range srcString {
		matched := true
		for j := range expectedSubstr {
			if srcString[i] != expectedSubstr[j] {
				matched = false
				break
			}
			i++
			if i >= len(srcString) {
				break
			}
		}
		if matched {
			if outCtr >= len(out) {
				return -1 // Need more space
			}
			out[outCtr] = i - len(expectedSubstr)
			outCtr++
		}
	}
	return outCtr
}
