package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func calculate(input string) (float64, error) {
	input = strings.TrimSpace(input)
	parts := strings.Split(input, " ")
	if len(parts) < 3 {
		return 0, fmt.Errorf("wrong format, parts less than 3")
	}

	num1, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("wrong format, %s is not a number", parts[0])
	}

	num2, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("wrong format, %s is not a number", parts[2])
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
			return 0, fmt.Errorf("division by zero")
		}
		return float64(num1) / float64(num2), nil
	default:
		return 0, fmt.Errorf("wrong format")
	}
}

func handleConnection(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	// Сообщаем WaitGroup, что горутина завершила работу при выходе из функции
	defer wg.Done()
	defer conn.Close()

	// 2. Создаем контекстный логгер для конкретного клиента.
	// Теперь КАЖДОЕ сообщение этого логгера будет автоматически содержать IP клиента!
	clientAddr := conn.RemoteAddr().String()
	logger := slog.With("client_addr", clientAddr)

	logger.Info("Client connected")

	// Локальный канал для сигнала о том, что обработка клиента завершена
	done := make(chan struct{})
	defer close(done) // Закроется автоматически при выходе из handleConnection

	// Безопасная фоновая горутина
	go func() {
		select {
		case <-ctx.Done():
			// Сервер останавливается -> принудительно рвем соединение
			_ = conn.Close()
		case <-done:
			// Клиент отключился сам -> просто выходим из горутины, без утечек!
			return
		}
	}()

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		req := strings.TrimSpace(scanner.Text())
		if req == "" {
			continue
		}

		result, err := calculate(req)
		var output string
		if err != nil {
			// Логируем ошибку формата с дополнительным полем "request"
			logger.Warn("Calculation failed", "request", req, "error", err.Error())
			output = fmt.Sprintf("Error: %v\n", err)
		} else {
			output = fmt.Sprintf("%f\n", result)
		}

		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Minute))

		if _, writeErr := conn.Write([]byte(output)); writeErr != nil {
			logger.Error("Failed to write response", "error", writeErr.Error())
			return
		}
	}

	if err := scanner.Err(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		logger.Error("Scanner error", "error", err.Error())
	} else {
		logger.Info("Client disconnected")
	}
}

// server теперь принимает контекст и адрес (порт) для запуска
func server(ctx context.Context, addr string) error {
	var lc net.ListenConfig
	// Использование ListenConfig позволяет привязать контекст к слушателю
	listener, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer listener.Close()

	log.Printf("Сервер калькулятора запущен на %s\n", addr)

	// WaitGroup для отслеживания всех активных клиентских соединений
	var wg sync.WaitGroup

	// Горутина для закрытия лисенера при отмене контекста
	go func() {
		<-ctx.Done()
		log.Println("Получен сигнал остановки, закрываем новые подключения...")
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			// Если ошибка вызвана закрытием лисенера — выходим из цикла штатно
			if errors.Is(err, net.ErrClosed) {
				break
			}
			log.Printf("Ошибка Accept: %v\n", err)
			continue
		}

		wg.Add(1)
		// Передаем контекст и WaitGroup в обработчик
		go handleConnection(ctx, conn, &wg)
	}

	log.Println("Ожидаем завершения обработки текущих клиентов...")
	wg.Wait() // Ждем, пока все активные клиенты дочитают/допишут данные
	log.Println("Все соединения закрыты. Сервер успешно остановлен.")
	return nil
}

func main() {
	// Убираем хардкод порта: теперь его можно задать через флаг -port=8080
	port := flag.String("port", "9999", "TCP port to listen on")
	flag.Parse()

	addr := net.JoinHostPort("", *port)

	// 1. Инициализируем JSON-логгер в самом начале программы
	// Вся программа теперь будет писать логи в красивом JSON-формате
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(jsonHandler))

	slog.Info("Starting calculator server...")

	// Создаем контекст, который отменится при системных сигналах прерывания
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server(ctx, addr); err != nil {
		log.Fatalf("Критическая ошибка сервера: %v\n", err)
	}
}
