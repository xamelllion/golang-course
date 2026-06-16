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

func (h *Handler) GetRepository(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return
	}

	params := strings.Split(r.URL.Path[1:], "/")
	if len(params) != 2 || params[0] == "" || params[1] == "" {
		http.Error(w, "invalid repository path", http.StatusBadRequest)
		return
	}

	owner, repo := params[0], params[1]

	repository, error := h.RepositoryUseCase.GetRepository(owner, repo)
	if error != nil {
		switch apperror.CodeOf(error) {
		case apperror.CodeNotFound:
			http.Error(w, "not found", http.StatusNotFound)
		case apperror.CodeInvalidArgument:
			http.Error(w, "invalid owner or repo", http.StatusBadRequest)
		case apperror.CodeUnavailable:
			http.Error(w, "unavailable github service", http.StatusBadGateway)
		default:
			fmt.Fprintf(os.Stderr, "internal error: %v\n", error)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(repository)
}
