package handler

import (
	"net/http"
	"strconv"
	"test-task/pkg/mid"
	"test-task/pkg/service"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handler struct {
	mos         *service.MosService
	Router      *chi.Mux
	middlewares []mid.Middleware
}

func (h *Handler) applyMiddlewares() {
	if h.middlewares == nil {
		return
	}
	for _, mw := range h.middlewares {
		h.Router.Use(mw)
	}
}

func (h *Handler) Mux() *chi.Mux {
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
			r.HandleFunc("/all", h.HandleAllParkings)
		})

	})

	return r
}

func NewHandler(mos *service.MosService, middlewares ...mid.Middleware) *Handler {
	h := &Handler{
		mos: mos,
	}

	h.Router = h.Mux()
	h.middlewares = append(h.middlewares, middlewares...)
	h.applyMiddlewares()

	return h
}

func (h *Handler) HandleAuth(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) HandleParkingByGlobalId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if !validateId(id) {
		render.Status(r, http.StatusBadRequest)
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

func (h *Handler) HandleAllParkings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.mos.GetRowsFromStorage(ctx)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		return
	}

	render.JSON(w, r, rows)
	render.Status(r, http.StatusOK)
}
func (h *Handler) HandleParkingById(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) HandleParkingByMode(w http.ResponseWriter, r *http.Request) {}

func validateId(id string) bool {
	intId, err := strconv.Atoi(id)
	if err != nil || id == "" {
		return false
	}
	return intId > 0
}
