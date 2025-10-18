package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	apiv1 "github.com/sharifahmad2061/trip-grpc-go/api/gen/go"
	"github.com/sharifahmad2061/trip-grpc-go/internal/db"
	queries "github.com/sharifahmad2061/trip-grpc-go/internal/db/generated"
	"github.com/sharifahmad2061/trip-grpc-go/internal/service"
	"github.com/sharifahmad2061/trip-grpc-go/internal/telemetry"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/log/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

func main() {
	// initialize telemetry here (tracing, metrics, etc.)
	ctx := context.Background()
	shutdown, err := telemetry.Init(ctx)
	if err != nil {
		panic(err)
	}
	defer shutdown()

	// logger
	lp := global.GetLoggerProvider()
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		otelzap.NewCore("trip-grpc-go", otelzap.WithLoggerProvider(lp)),
	)
	logger := zap.New(core)
	zap.ReplaceGlobals(logger)
	grpclog.SetLoggerV2(zapgrpc.NewLogger(logger))
	logger.Info("Logger initialized")

	// db connection
	db, err := db.Initialize(ctx)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
		panic(err)
	}
	defer db.Close()

	// network socket
	socket, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to listen: %v", err))
	}

	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	query := queries.New(db)
	apiv1.RegisterTripsServer(server, &service.TripsServiceImpl{Query: query})
	reflection.Register(server)

	// start runtime telemetry (memory, GC, etc.) collection
	_ = runtime.Start()

	go func() {
		logger.Info("gRPC server is running on port 50051")
		if err := server.Serve(socket); err != nil {
			logger.Fatal(fmt.Sprintf("failed to serve: %v", err))
		}
		logger.Info("Server stopped serving")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")
	server.GracefulStop()
	logger.Info("Server stopped")
}
