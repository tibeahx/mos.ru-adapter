package handler

import (
	"fmt"
	"net/http"

	"github.com/tibeahx/mos.ru-adapter/pkg/helper"
	"github.com/tibeahx/mos.ru-adapter/pkg/mid"
	"github.com/tibeahx/mos.ru-adapter/pkg/svc/mos"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	globalId  = "globalId"
	id        = "id"
	parkingId = "parkingId"
	mode      = "mode"
	jsonErr   = "error"
	result    = "result"
)

type Handler struct {
	mos         *mos.Mos
	Mux         *chi.Mux
	middlewares []mid.Middleware
}

func (h *Handler) applyMiddlewares() {
	if h.middlewares == nil {
		return
	}
	for _, mw := range h.middlewares {
		h.Mux.Use(mw)
	}
}

func (h *Handler) initMux() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/find/parkings", func(r chi.Router) {
			r.HandleFunc("/globalId", h.HandleParkingByGlobalId)
			r.HandleFunc("/id", h.HandleParkingById)
			r.HandleFunc("/mode", h.HandleParkingByMode)
		})

	})

	return r
}

func NewHandler(mos *mos.Mos, middlewares ...mid.Middleware) *Handler {
	h := &Handler{
		mos: mos,
	}

	h.Mux = h.initMux()
	h.middlewares = append(h.middlewares, middlewares...)
	h.applyMiddlewares()

	return h
}

func (h *Handler) HandleAuth(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) HandleParkingByGlobalId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "globalId")

	if !helper.VadlidateID(id) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{jsonErr: fmt.Sprintf("%s: %s not found", globalId, id)})
		return
	}

	parking, err := h.mos.GetParkingByGlobalId(ctx, id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, parking)
	render.Status(r, http.StatusOK)
}
func (h *Handler) HandleParkingById(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) HandleParkingByMode(w http.ResponseWriter, r *http.Request) {}
