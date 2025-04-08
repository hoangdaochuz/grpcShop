package adapter

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"grpcShop.com/backend/apis/product"
	"grpcShop.com/backend/config"
	"grpcShop.com/backend/db"
	productService "grpcShop.com/backend/product"
)

type ProductService struct {
	config       *config.Config
	grpcServer   *grpc.Server
	httpServer   *http.Server
	grpcListener net.Listener
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewProductService(cfg *config.Config) *ProductService {
	ctx, cancel := context.WithCancel(context.Background())
	return &ProductService{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}

func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func (prodService *ProductService) Start() error {
	db := db.NewDB()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	prodService.grpcListener = lis
	prodService.grpcServer = grpc.NewServer()

	product.RegisterProductServiceServer(prodService.grpcServer, &productService.ProductServer{
		Db: db})
	// Khởi chạy gRPC server trong goroutine
	prodService.wg.Add(1)
	go func() {
		fmt.Println("Starting gRPC server on :50051")
		defer prodService.wg.Done()
		if err := prodService.grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// Khởi tạo gRPC Gateway
	mux := runtime.NewServeMux()
	err = product.RegisterProductServiceHandlerServer(context.Background(), mux, &productService.ProductServer{
		Db: db})
	if err != nil {
		return err
	}

	// Start HTTP server
	prodService.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: corsHandler(mux),
	}

	prodService.wg.Add(1)
	go func() {
		defer prodService.wg.Done()
		fmt.Println("gRPC-Gateway is running on :8080")
		err = prodService.httpServer.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

func (prodService *ProductService) Stop() error {
	// Dừng HTTP server
	if prodService.httpServer != nil {
		if err := prodService.httpServer.Shutdown(prodService.ctx); err != nil {
			return err
		}
	}

	if prodService.grpcServer != nil {
		prodService.grpcServer.GracefulStop()
	}

	// Đóng listener
	if prodService.grpcListener != nil {
		prodService.grpcListener.Close()
	}

	// Hủy context
	prodService.cancel()
	// Chờ cho tất cả goroutine hoàn thành
	prodService.wg.Wait()
	return nil
}
