package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/broswen/kvs/internal/cache"
	"github.com/broswen/kvs/internal/db"
	"github.com/broswen/kvs/internal/handlers"
	"github.com/broswen/kvs/internal/item"
	"github.com/broswen/kvs/pkg/kvs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server interface {
	Start() error
}

type ChiServer struct {
	cacheService cache.CacheService
	dbService    db.DBService
	itemService  item.ItemService
	logger       zerolog.Logger
	router       chi.Router
}

type GRPCServer struct {
	kvs.UnimplementedKeyValueServiceServer
	cacheService cache.CacheService
	dbService    db.DBService
	itemService  item.ItemService
	logger       zerolog.Logger
	router       chi.Router
}

func New() (ChiServer, error) {
	cacheService, err := cache.New()
	if err != nil {
		return ChiServer{}, fmt.Errorf("init cacheService: %w\n", err)
	}
	dbService, err := db.New()
	if err != nil {
		return ChiServer{}, fmt.Errorf("init dbService: %w\n", err)
	}

	itemService, err := item.New(cacheService, dbService)
	if err != nil {
		return ChiServer{}, fmt.Errorf("init itemService: %w\n", err)
	}

	logger := httplog.NewLogger("kvs", httplog.Options{
		JSON: true,
	})
	server := ChiServer{
		itemService: itemService,
		logger:      logger,
		router:      chi.NewRouter(),
	}
	server.SetRoutes()
	return server, nil
}

func (s ChiServer) Start() error {
	s.logger.Info().Msg("Starting chi server...")
	return http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), s.router)
}

func (s ChiServer) SetRoutes() {
	s.router.Use(httplog.RequestLogger(s.logger))
	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	// health check
	s.router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	s.router.Get("/{key}", handlers.GetHandler(s.itemService))
	s.router.Post("/{key}", handlers.SetHandler(s.itemService))
	s.router.Delete("/{key}", handlers.DeleteHandler(s.itemService))
}

func NewGRPC() (Server, error) {
	cacheService, err := cache.New()
	if err != nil {
		return ChiServer{}, fmt.Errorf("init cacheService: %w\n", err)
	}
	dbService, err := db.New()
	if err != nil {
		return ChiServer{}, fmt.Errorf("init dbService: %w\n", err)
	}

	itemService, err := item.New(cacheService, dbService)
	if err != nil {
		return ChiServer{}, fmt.Errorf("init itemService: %w\n", err)
	}

	logger := httplog.NewLogger("kvs", httplog.Options{
		JSON: true,
	})
	server := GRPCServer{
		itemService: itemService,
		logger:      logger,
		router:      chi.NewRouter(),
	}
	return server, nil
}

func (s GRPCServer) Start() error {
	s.logger.Info().Msg("Starting gRPC server...")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", os.Getenv("PORT")))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	kvs.RegisterKeyValueServiceServer(grpcServer, s)
	return grpcServer.Serve(lis)
}

func (s GRPCServer) SetValue(ctx context.Context, req *kvs.SetValueRequest) (*kvs.SetValueResponse, error) {
	err := s.itemService.Set(item.Item{Key: req.Key, Value: req.Value})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}
	return &kvs.SetValueResponse{
		Key: req.Key,
	}, nil
}

func (s GRPCServer) GetValue(ctx context.Context, req *kvs.GetValueRequest) (*kvs.GetValueResponse, error) {
	foundItem, err := s.itemService.Get(req.Key)
	if err != nil {
		if errors.As(err, &item.ErrItemNotFound{}) {
			return nil, status.Errorf(codes.NotFound, "Item not found: %s", req.Key)
		}
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}
	return &kvs.GetValueResponse{
		Key:   foundItem.Key,
		Value: foundItem.Value,
	}, nil
}

func (s GRPCServer) DeleteValue(ctx context.Context, req *kvs.DeleteValueRequest) (*kvs.DeleteValueResponse, error) {
	err := s.itemService.Delete(req.Key)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}
	return &kvs.DeleteValueResponse{
		Key: req.Key,
	}, nil
}
