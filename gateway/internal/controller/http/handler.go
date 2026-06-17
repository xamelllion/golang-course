package controller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "github.com/xamelllion/golang-course/docs"
	"github.com/xamelllion/golang-course/gateway/internal/domain"
	apperror "github.com/xamelllion/golang-course/internal/errors"
)

type RepositoryUseCase interface {
	GetRepository(owner, repo string) (domain.Repository, error)
}

type PingUseCase interface {
	PingAll() (domain.ServicesInfo, error)
}

type Handler struct {
	RepositoryUseCase RepositoryUseCase
	PingUseCase       PingUseCase
	log               *slog.Logger
}

func NewHandler(repositoryUsecase RepositoryUseCase, pingUsecase PingUseCase, log *slog.Logger) Handler {
	return Handler{
		RepositoryUseCase: repositoryUsecase,
		PingUseCase:       pingUsecase,
		log:               log,
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// GetRepository godoc
//
//	@Summary		Get repository by GitHub URL
//	@Description	get repository information by GitHub repository URL
//	@ID				get-repo-by-owner-repo
//	@Tags			repo
//	@Accept			json
//	@Produce		json
//	@Param			url	query		string	true	"GitHub repository URL"
//	@Success		200		{object}	domain.Repository
//	@Failure		400		{object}	controller.GetRepository.HTTPError
//	@Failure		404		{object}	controller.GetRepository.HTTPError
//	@Failure		502		{object}	controller.GetRepository.HTTPError
//	@Failure		500		{object}	controller.GetRepository.HTTPError
//	@Router			/api/repositories/info [get]
func (h *Handler) GetRepository(w http.ResponseWriter, r *http.Request) {
	type HTTPError struct {
		Err string `json:"error" example:"some error"`
	}

	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "wrong method")
		return
	}

	githubURLstring := r.URL.Query().Get("url")
	if githubURLstring == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid repository path")
		return
	}

	githubURL, err := url.Parse(githubURLstring)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid repository path")
		return
	}

	if githubURL.Host != "github.com" {
		writeJSONError(w, http.StatusBadRequest, "invalid repository path")
		return
	}

	if githubURL.Path == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid repository path")
		return
	}

	if githubURL.Path[0] == '/' {
		githubURL.Path = githubURL.Path[1:]
	}

	params := strings.Split(githubURL.Path, "/")
	if len(params) < 2 {
		writeJSONError(w, http.StatusBadRequest, "invalid repository path")
		return
	}

	owner := params[0]
	repo := params[1]

	if owner == "" || repo == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid repository path")
		return
	}

	h.log.Debug("http: get request", "owner", owner, "repo", repo)

	repository, error := h.RepositoryUseCase.GetRepository(owner, repo)
	if error != nil {
		h.log.Debug("http: error", "error", error)
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

// PingServices godoc
//
//	@Summary		Get status of services
//	@Description	get status of services
//	@ID				get-status-of-services
//	@Tags			service
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	domain.ServicesInfo
//	@Failure		503		{object}	domain.ServicesInfo
//	@Failure		500		{object}	controller.GetRepository.HTTPError
//	@Router			/api/ping [get]
func (h *Handler) PingServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "wrong method")
		return
	}
	h.log.Debug("http: ping request")

	servicesInfo, err := h.PingUseCase.PingAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if servicesInfo.Status == domain.ServicesStatusOk {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	_ = json.NewEncoder(w).Encode(servicesInfo)
}
