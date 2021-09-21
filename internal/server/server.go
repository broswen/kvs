package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/broswen/kvs/internal/cache"
	"github.com/broswen/kvs/internal/db"
	"github.com/broswen/kvs/internal/handlers"
	"github.com/broswen/kvs/internal/item"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

type Server struct {
	cacheService cache.CacheService
	dbService    db.DBService
	itemService  item.ItemService
	logger       zerolog.Logger
	router       chi.Router
}

func New() (Server, error) {
	cacheService, err := cache.New()
	if err != nil {
		return Server{}, fmt.Errorf("init cacheService: %w\n", err)
	}
	dbService, err := db.New()
	if err != nil {
		return Server{}, fmt.Errorf("init dbService: %w\n", err)
	}

	itemService, err := item.New(cacheService, dbService)
	if err != nil {
		return Server{}, fmt.Errorf("init itemService: %w\n", err)
	}

	logger := httplog.NewLogger("kvs", httplog.Options{
		JSON: true,
	})
	return Server{
		itemService: itemService,
		logger:      logger,
		router:      chi.NewRouter(),
	}, nil
}

func (s Server) Start() error {
	s.logger.Info().Msg("Starting server...")
	return http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), s.router)
}

func (s Server) SetRoutes() {
	s.router.Use(httplog.RequestLogger(s.logger))
	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	// health check
	s.router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	s.router.Get("/{key}", handlers.GetHandler(s.itemService))
	s.router.Post("/{key}", handlers.SetHandler(s.itemService))
}
