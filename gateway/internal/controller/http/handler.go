package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/xamelllion/golang-course/gateway/internal/domain"
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
		fmt.Fprintf(os.Stderr, "error: %s", error.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(repository)
}
