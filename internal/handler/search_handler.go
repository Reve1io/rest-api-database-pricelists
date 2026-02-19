package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"rest-api-database-pricelists/internal/dto"
	"rest-api-database-pricelists/internal/service"

	"go.uber.org/zap"
)

type SearchHandler struct {
	service *service.SearchService
	logger  *zap.Logger
}

func NewSearchHandler(s *service.SearchService, l *zap.Logger) *SearchHandler {
	return &SearchHandler{service: s, logger: l}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {

	h.logger.Info("incoming request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	var req dto.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Data) < 2 {
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}

	row := req.Data[1]

	mpn, ok := row[0].(string)
	if !ok {
		http.Error(w, "invalid mpn format", http.StatusBadRequest)
		return
	}

	var qty int

	switch v := row[1].(type) {
	case float64:
		qty = int(v)
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, "invalid quantity format", http.StatusBadRequest)
			return
		}
		qty = parsed
	default:
		http.Error(w, "invalid quantity type", http.StatusBadRequest)
		return
	}

	result, err := h.service.Search(r.Context(), mpn, qty)
	if err != nil {
		h.logger.Error("search failed", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

	h.logger.Info("sending response",
		zap.Int("items_count", len(result)),
	)
}
