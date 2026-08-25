// пример как код был реализован без контекстов
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func calculate(input string) (float64, error) {
	parts := strings.Split(input, " ")

	if len(parts) < 3 {
		return 0, fmt.Errorf("wrong format, parts less then 3")
	}

	num1, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("wrong format, %s not number", parts[0])
	}

	num2, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("wrong format, %s not number", parts[2])
	}

	if len(parts[1]) != 1 || !strings.Contains("+-*/:", parts[1]) {
		return 0, fmt.Errorf("wrong format, %s not in '+-*/:'", parts[1])
	}

	switch parts[1] {
	case "+":
		return float64(num1 + num2), nil
	case "-":
		return float64(num1 - num2), nil
	case "*":
		return float64(num1 * num2), nil
	case "/", ":":
		if num2 == 0 {
			return 0, fmt.Errorf("divide by zero")
		}
		return float64(num1) / float64(num2), nil
	default:
		return 0, fmt.Errorf("wrong format")
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Защита от вечных соединений (например, 2 минуты на операцию)
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Minute))

	// Идеально для TCP: читает поток построчно до \n
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		req := strings.TrimSpace(scanner.Text())
		if req == "" {
			continue
		}
		log.Printf("Прочитано %d байт: %s\n", len(req), req)
		result, err := calculate(req)
		var output string
		if err != nil {
			output = fmt.Sprintf("Error: %v\n", err)
		} else {
			output = fmt.Sprintf("%.2f\n", result)
		}
		fmt.Println(output)

		// Обновляем таймаут для следующей команды
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Minute))

		if _, writeErr := conn.Write([]byte(output)); writeErr != nil {
			log.Printf("Ошибка записи: %v\n", writeErr)
			return
		}

	}

	if err := scanner.Err(); err != nil {
		// Логируем реальную ошибку чтения, если она была
		log.Printf("Ошибка чтения: %v\n", err)
	}
}

func server(port int) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка Accept %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func main() {
	fmt.Println("Запущен сервер калькулятора, пришлите запись формата 1 + 2")

	port := flag.Int("port", 9999, "Run as client")
	flag.Parse()
	server(*port)
}
