package main

import (
//	"fmt"
)

func countTrailingBits(mask uint64) uint64 {
	mask = ^mask
	for i := uint64(0); i < 64; i++ {
		if mask&(1<<i) != 0 {
			return uint64(i)
		}
	}
	return 64
}

func popCount(in uint64) uint8 {
	var ret uint8
	var i uint8
	for i = 0; i < 64; i++ {
		if in&(1<<i) != 0 {
			ret++
		}
	}
	return ret
}

func permuteInt(orig, mask uint64) uint64 {
	var origCnt uint8
	var maskCnt uint8
	var ret uint64

	for {
		if maskCnt >= 64 {
			break
		}
		if origCnt >= 64 {
			break
		}

		//		fmt.Println(origCnt, maskCnt)
		if orig&(1<<origCnt) != 0 {
			ret |= 1 << maskCnt
		}

		origCnt++
		maskCnt++
		for {
			if mask&(1<<maskCnt) != 0 {
				break
			} else {
				maskCnt++
			}
			if maskCnt >= 64 {
				break
			}
		}
	}
	return ret
}
