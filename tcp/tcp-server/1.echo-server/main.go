package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func handleConnection(c net.Conn) {
	defer c.Close() // Гарантированно закрываем сокет на выходе

	buf := make([]byte, 1024) // Буфер побольше для оптимизации сисколов
	for {
		l, err := c.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Клиент успешно закрыл соединение")
				break
			}
			log.Printf("Ошибка чтения: %v\n", err)
			return // Выходим из горутины, не роняя сервер
		}

		fmt.Printf("Прочитано %d байт: %s\n", l, buf[:l])

		// Пишем обратно ровно столько, сколько прочитали
		_, writeErr := c.Write(buf[:l])
		if writeErr != nil {
			log.Printf("Ошибка записи: %v\n", writeErr)
			return
		}
	}
}

func server() {
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		// Wait for a connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка Accept: %v\n", err)
			continue // Продолжаем слушать другие подключения
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handleConnection(conn)
	}
}

func client() {
	// 1. Устанавливаем соединение с таймаутом на подключение
	dialer := net.Dialer{Timeout: 5 * time.Second}
	conn, err := dialer.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Printf("Ошибка подключения: %v", err)
		return
	}
	defer conn.Close()

	// 2. Устанавливаем общий таймаут на сетевые операции
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// 3. Отправляем данные с обязательной проверкой ошибки
	message := []byte("message from client")
	_, err = conn.Write(message)
	if err != nil {
		log.Printf("Ошибка записи: %v", err)
		return
	}

	// 4. Безопасно читаем ВЕСЬ ответ до закрытия соединения сервером
	// response, err := io.ReadAll(conn) не подходит для этого примера, сервер не закрывает соединение
	
	// ЧИТАЕМ ОТВЕТ БЕЗ io.ReadAll
	buf := make([]byte, 1024)
	l, err := conn.Read(buf) // Прочитает ровно то эхо, которое сервер только что вернул
	if err != nil && err != io.EOF {
		log.Printf("Ошибка чтения: %v", err)
		return
	}

	fmt.Printf("Прочитано %d байт: %s\n", l, buf[:l])
}

func main() {
	isClient := flag.Bool("client", false, "Run as client")
	isServer := flag.Bool("server", false, "Run as client")
	flag.Parse()

	if *isServer {
		fmt.Println("[*] Running server")
		server()
		return
	}

	if *isClient {
		client()
		return
	}
}
