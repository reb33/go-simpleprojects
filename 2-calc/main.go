package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func main() {
	operation := scanOperation()
	arr := scanArray()
	res := calculate(operation, arr)
	fmt.Printf("Результат операции %v: %.2f\n", operation, res)
}

func scanOperation() string {
	var operation string
	for {
		fmt.Println("Введите операцию: AVG, SUM, MED")
		fmt.Scan(&operation)
		if operation != "AVG" && operation != "avg" && operation != "SUM" && operation != "sum" &&
			operation != "MED" && operation != "med" {
			fmt.Println("Некорректная операция")
			continue
		}
		break
	}
	return strings.ToUpper(operation) // Convert to uppercase for consistency in 
}

func scanArray() []int {
	var arr []int
	ans := "n"
	for len(arr) < 2 && (ans != "y" && ans != "Y"){
		fmt.Println("Введите массив чисел через запятую")
		// fmt.Scan(&arrString)  // не подходит Scan, потому что при прробеле он прерывает ввод
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		splitedArr := strings.SplitSeq(input, ",")
		for v := range splitedArr {
			trimedStr := strings.TrimSpace(v)
			num, error := strconv.Atoi(trimedStr)
			if error != nil {
				continue
			}
			arr = append(arr, num)
		}
		fmt.Printf("Получился массив: %v\nПодтвердить (y/n): ", arr)
		fmt.Scan(&ans)
	}
	return arr
}

func calculate(operation string, arr []int) float64 {
	switch operation {
		case "AVG":
			return float64(sum(arr)) / float64(len(arr))
		case "SUM":
			return float64(sum(arr))
		case "MED":
			return med(arr)
	}
	return 0
}

func sum(arr []int) int {
	var sum int
	for _, v := range arr {
		sum += v
	}
	return sum
}

func med(arr []int) float64 {
	arr_copy := slices.Clone(arr)
	sort.Ints(arr_copy)
	mid := len(arr_copy) / 2

	if len(arr_copy) % 2 == 1 {
		return float64(arr_copy[mid])
	}else {
		return float64(arr_copy[mid-1]+arr_copy[mid]) / 2
	}
}