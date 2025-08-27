package main

import (
	"fmt"
	"math/big"
)


func faktorial(n int) *big.Int {
	result := big.NewInt(1)
	for i := 2; i <= n; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}
	return result
}


func f(n int) *big.Int {
	fact := faktorial(n)

	power := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(n)), nil)


	quotient := new(big.Int).Div(fact, power)

	remainder := new(big.Int).Mod(fact, power)
	if remainder.Cmp(big.NewInt(0)) > 0 {
		quotient.Add(quotient, big.NewInt(1)) // pembulatan ke atas
	}

	return quotient
}

func main() {
var n int
	fmt.Print("Masukkan nilai angka yang ingin difaktor: ")
	fmt.Scan(&n)

	result := f(n)
	fmt.Printf("f(%d) = %v\n", n, result)
}
