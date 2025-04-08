package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"grpcShop.com/backend/adapter"
	"grpcShop.com/backend/config"
)

type Service interface {
	Start() error
	Stop() error
}

func initServices(cfg *config.Config) []Service {
	// Khởi tạo các service
	productService := adapter.NewProductService(cfg)
	return []Service{
		productService,
	}
}

func runService(service Service, idx int, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	err := service.Start()
	if err != nil {
		errChan <- err
		return
	}
}

func main() {
	// Load config
	cfg := config.LoadConfig()
	var wg sync.WaitGroup

	// Channel để theo dõi lỗi từ các service
	errChan := make(chan error, 10)

	// Channel để xử lý tín hiệu kết thúc
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	services := initServices(cfg)

	// Khởi chạy tất cả service
	for i, service := range services {
		wg.Add(1)
		go runService(service, i, &wg, errChan)
	}

	for {

		select {
		case err := <-errChan:
			fmt.Printf("Error from service: %v\n", err)
		case sig := <-sigChan:

			fmt.Printf("Received signal: %v\n", sig)
			// Dừng tất cả service
			for _, service := range services {
				if err := service.Stop(); err != nil {
					fmt.Printf("Error stopping service: %v\n", err)
				}
			}
			wg.Wait()
		}
	}

}
