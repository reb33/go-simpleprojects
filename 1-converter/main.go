package main

import (
	"fmt"
	"strings"
)

const usdToEur = 0.85
const usdToRub = 75.0
const eurToRub = usdToRub / usdToEur

func main() {
	currFrom, currTo, valueToConvert := scanData()
	convertedValue := convertCurrency(valueToConvert, currFrom, currTo)
	fmt.Printf("Конвертация из %v в %v значения: %v получилось: %.2f\n", currFrom, currTo, valueToConvert, convertedValue)
}

func scanData() (string, string, float64) {
	currFrom := scanCurrFrom()
	valueToConvert := scanValueToConvert()
	currTo := scanCurrTo(currFrom)
	return currFrom, currTo, valueToConvert
}

func scanCurrFrom() (string) {
	var currFrom string
	for {
		fmt.Println("Введите валюту, из которой хотите конвертировать: USD, RUB, EUR")
		fmt.Scan(&currFrom)
		if currFrom != "USD" && currFrom != "usd" && currFrom != "EUR" && currFrom != "eur" && currFrom != "RUB" && currFrom != "rub" {
			fmt.Println("Некорректная валюта")
		}else{
			break
		}
	}
	return strings.ToUpper(currFrom)
}

func scanValueToConvert() (float64) {
	var valueToConvert float64
	for {
		fmt.Println("Введите значение для конвертации")
		fmt.Scan(&valueToConvert)
		if valueToConvert <= 0 {
			fmt.Println("Значение должно быть положительным числом")
			continue
		}
		break
	}
	return valueToConvert
}

func scanCurrTo(currFrom string) (string) {
	var currTo string
	outerLoop: for {
		switch  currFrom{
		case "USD":
			fmt.Println("Введите валюту, в которую хотите конвертировать: RUB, EUR")
			fmt.Scan(&currTo)
			if currTo != "RUB" && currTo != "rub" && currTo != "EUR" && currTo != "eur" {
				fmt.Println("Некорректная валюта")
				continue
			}
			break outerLoop
		case "EUR":
			fmt.Println("Введите валюту, в которую хотите конвертировать: USD, RUB")
			fmt.Scan(&currTo)
			if currTo != "USD" && currTo != "usd" && currTo != "RUB" && currTo != "rub" {
				fmt.Println("Некорректная валюта")
				continue
			}
			break outerLoop
		case "RUB":
			fmt.Println("Введите валюту, в которую хотите конвертировать: USD, EUR")
			fmt.Scan(&currTo)
			if currTo != "USD" && currTo != "usd" && currTo != "EUR" && currTo != "eur" {
				fmt.Println("Некорректная валюта")
				continue
			}
			break outerLoop
		}
	}
	return strings.ToUpper(currTo)	
}

func convertCurrency(amount float64, from string, to string) float64 {
	switch from{
	case "USD":
		switch to{
			case "RUB":
				return amount * usdToRub
			case "EUR":
				return amount * usdToEur
		}
	case "EUR":
		switch to{
			case "RUB":
				return amount * eurToRub
			case "USD":
				return amount / usdToEur
		}
	case "RUB":
		switch to{
			case "USD":
				return amount / usdToRub
			case "EUR":
				return amount / eurToRub
		}
	}
	return 0
}