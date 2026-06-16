package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/xamelllion/golang-course/gateway/internal/domain"
	apperror "github.com/xamelllion/golang-course/internal/errors"
)

type RepositoryUseCase interface {
	GetRepository(owner, repo string) (domain.Repository, error)
}

type Handler struct {
	RepositoryUseCase RepositoryUseCase
}

func NewHandler(rp RepositoryUseCase) Handler {
	return Handler{
		RepositoryUseCase: rp,
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func (h *Handler) GetRepository(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "wrong method")
		return
	}

	params := strings.Split(r.URL.Path[1:], "/")
	if len(params) != 2 || params[0] == "" || params[1] == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid repository path")
		return
	}

	owner, repo := params[0], params[1]

	repository, error := h.RepositoryUseCase.GetRepository(owner, repo)
	if error != nil {
		switch apperror.CodeOf(error) {
		case apperror.CodeNotFound:
			writeJSONError(w, http.StatusNotFound, "not found")
		case apperror.CodeInvalidArgument:
			writeJSONError(w, http.StatusBadRequest, "invalid owner or repo")
		case apperror.CodeUnavailable:
			writeJSONError(w, http.StatusBadGateway, "unavailable github service")
		default:
			fmt.Fprintf(os.Stderr, "internal error: %v\n", error)
			writeJSONError(w, http.StatusInternalServerError, "internal error")
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(repository)
}
