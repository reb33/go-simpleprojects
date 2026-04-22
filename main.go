package main

import "fmt"


func main() {
	const usdToEur = 0.85
	const usdToRub = 75.0
	const eurToRub = usdToRub / usdToEur
}

func scanData() (string) {
	var data string
	fmt.Scan(&data)
	return data
}

func convertCurrency(amount float64, from string, to string) {}