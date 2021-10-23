package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/broswen/kvs/internal/item"
	items "github.com/broswen/kvs/internal/item"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
)

type GetResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (g *GetResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func GetHandler(itemService items.ItemGetterSetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		var key string
		if key = chi.URLParam(r, "key"); key == "" {
			oplog.Warn().Msg("must specify key")
			render.Render(w, r, ErrInvalidRequest(errors.New("must specify key")))
			return
		}

		item, err := itemService.Get(key)
		if err != nil {
			if errors.As(err, &items.ErrItemNotFound{}) {
				oplog.Warn().Msgf("item not found")
				render.Render(w, r, ErrNotFound(errors.New("item not found")))
				return
			}
			oplog.Err(err).Msg("error getting item")
			render.Render(w, r, ErrInternalServer(err))
			return
		}

		_, err = w.Write([]byte(item.Value))
		if err != nil {
			oplog.Err(err).Msg("error getting item")
			render.Render(w, r, ErrInternalServer(err))
			return
		}

	}
}

type SetResponse struct {
}

func (s *SetResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func SetHandler(itemService item.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		var key string
		if key = chi.URLParam(r, "key"); key == "" {
			oplog.Warn().Msg("must specify key")
			render.Render(w, r, ErrInvalidRequest(errors.New("must specify key")))
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			oplog.Err(err).Msg("error reading body")
			render.Render(w, r, ErrInternalServer(errors.New("error reading body")))
			return
		}
		defer r.Body.Close()

		item := item.Item{Key: key, Value: string(body)}
		err = itemService.Set(item)
		if err != nil {
			oplog.Err(err).Msg("error setting item")
			render.Render(w, r, ErrInternalServer(err))
			return
		}

		render.Render(w, r, &SetResponse{})
	}
}

type DeleteResponse struct {
}

func (s *DeleteResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func DeleteHandler(itemService item.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		var key string
		if key = chi.URLParam(r, "key"); key == "" {
			oplog.Warn().Msg("must specify key")
			render.Render(w, r, ErrInvalidRequest(errors.New("must specify key")))
			return
		}

		err := itemService.Delete(key)
		if err != nil {
			oplog.Err(err).Msg("error deleting item")
			render.Render(w, r, ErrInternalServer(err))
			return
		}

		render.Render(w, r, &SetResponse{})
	}
}
