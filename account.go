package main

import (
	"fmt"
	"math/big"
)

var (
	zeros = "000000000000000000"
)

func toWei(balance string) string {
	return fmt.Sprintf("%s%s", balance, zeros)
}

func toWeiHex(balance string) string {
	value, _ := new(big.Int).SetString(fmt.Sprintf("%s%s", balance, zeros), 10)
	return "0x" + value.Text(16)
}

func pkPadding(pk string) string {
	for len(pk) < 64 {
		pk = "0" + pk
	}
	return "0x" + pk
}
