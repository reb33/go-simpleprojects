package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func handlerConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	// Используем замыкание в defer, чтобы гарантировать, что закрытие произойдет,
	// даже если внутри функции будут паники (хорошая практика)
	defer func() { _ = conn.Close() }()

	clientAddr := conn.RemoteAddr().String()
	logger := slog.With("client_addr", clientAddr)

	logger.Debug("Client connected")

	// Устанавливаем таймаут на запись. Фоновая горутина с ctx здесь не нужна,
	// так как эта операция атомарна и мгновенна.
	_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	// Добавлен перенос строки \n в конце для корректного чтения клиентами
	response := time.Now().Format("02 Jan 2006 15:04:05\n")
	if _, writeError := conn.Write([]byte(response)); writeError != nil {
		logger.Error("Failed to write response", "error", writeError)
		return
	}
}

func server(ctx context.Context, address string) error {
	var lc net.ListenConfig

	listener, err := lc.Listen(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer listener.Close()

	slog.Info("Server started", "address", address)

	var wg sync.WaitGroup

	// Горутина для остановки listener при отмене контекста
	go func() {
		<-ctx.Done()
		slog.Warn("Got stop signal, closing listener...")
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break // Корректный выход при остановке сервера
			}
			slog.Error("Accept error", "error", err) // Исправлен формат логирования
			continue
		}

		wg.Add(1)
		// Передавать ctx в handler больше нет необходимости
		go handlerConnection(conn, &wg)
	}

	slog.Info("Waiting for active connections to change...")
	wg.Wait()
	slog.Info("All connections closed. Server stopped")
	return nil
}

func main() {
	port := flag.String("port", "9999", "TCP port to listen on")
	flag.Parse()

	addr := net.JoinHostPort("", *port)

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(jsonHandler))

	slog.Info("Starting time server...")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server(ctx, addr); err != nil {
		slog.Error("Critical server error", "error", err) // Исправлен log.Fatal с %w
		os.Exit(1)
	}
}