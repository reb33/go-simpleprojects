package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func handlerConnection(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	logger := slog.With("clent_addr", clientAddr)

	logger.Info("Client connected")

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			_ = conn.Close()
		case <-done:
			return
		}
	}()

	_ = conn.SetWriteDeadline(time.Now().Add(15 * time.Second))

	if _, writeError := conn.Write([]byte(time.Now().Format("02 Jan 2006 15:04:05"))); writeError != nil {
		logger.Error("Failed to write response", "error", writeError.Error())
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

	slog.Info("Server start on " + address)

	var wg sync.WaitGroup

	go func() {
		<-ctx.Done()
		slog.Warn("Got stop signal")
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			slog.Info("Ошибка Accept: %w", "error", err)
			continue
		}

		wg.Add(1)
		go handlerConnection(ctx, conn, &wg)
	}

	slog.Info("Ожидаем завершения обработки текущих клиентов...")
	wg.Wait()
	slog.Info("Все соединения закрыты. Сервер остановлен")
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
		log.Fatal("Критическая ошибка сервера: %w\n", err)
	}
}
