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
	GetSubscribedRepository() ([]domain.Repository, error)
}

type PingUseCase interface {
	PingAll() (domain.ServicesInfo, error)
}

type SubcribeUseCase interface {
	Subscribe(owner, repo string) error
	Unsubscribe(owner, repo string) error
	GetSubscriptions() ([]domain.Subscription, error)
}

type Handler struct {
	RepositoryUseCase RepositoryUseCase
	PingUseCase       PingUseCase
	SubscribeUseCase  SubcribeUseCase
	log               *slog.Logger
}

func NewHandler(repositoryUsecase RepositoryUseCase, pingUsecase PingUseCase, subscribeUsecase SubcribeUseCase, log *slog.Logger) Handler {
	return Handler{
		RepositoryUseCase: repositoryUsecase,
		PingUseCase:       pingUsecase,
		SubscribeUseCase:  subscribeUsecase,
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

type SubscribeRequest struct {
	Owner string `json:"owner" example:"octocat"`
	Repo  string `json:"repo" example:"Hello-World"`
}

// Subscribe godoc
//
//	@Summary		Subscribe to a repository
//	@Description	Create a subscription to an existing GitHub repository
//	@ID				subscribe-to-repository
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		SubscribeRequest	true	"Repository subscription"
//	@Success		200
//	@Failure		400	{object}	controller.GetRepository.HTTPError
//	@Failure		404	{object}	controller.GetRepository.HTTPError
//	@Failure		409	{object}	controller.GetRepository.HTTPError
//	@Failure		500	{object}	controller.GetRepository.HTTPError
//	@Failure		502	{object}	controller.GetRepository.HTTPError
//	@Router			/subscriptions [post]
func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	var subReq SubscribeRequest

	if err := json.NewDecoder(r.Body).Decode(&subReq); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.SubscribeUseCase.Subscribe(subReq.Owner, subReq.Repo); err != nil {
		h.log.Debug("http: error", "error", err)
		switch apperror.CodeOf(err) {
		case apperror.CodeNotFound:
			writeJSONError(w, http.StatusNotFound, "not found")
		case apperror.CodeInvalidArgument:
			writeJSONError(w, http.StatusBadRequest, "invalid owner or repo")
		case apperror.CodeUnavailable:
			writeJSONError(w, http.StatusBadGateway, "unavailable github service")
		case apperror.CodeDublicate:
			writeJSONError(w, http.StatusConflict, "already subscribed")
		default:
			fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
			writeJSONError(w, http.StatusInternalServerError, "internal error")
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}

// Unsubscribe godoc
//
//	@Summary		Unsubscribe from a repository
//	@Description	Delete an existing repository subscription
//	@ID				unsubscribe-from-repository
//	@Tags			subscriptions
//	@Produce		json
//	@Param			owner	path	string	true	"GitHub repository owner"	example(octocat)
//	@Param			repo	path	string	true	"GitHub repository name"		example(Hello-World)
//	@Success		200
//	@Failure		400	{object}	controller.GetRepository.HTTPError
//	@Failure		404	{object}	controller.GetRepository.HTTPError
//	@Failure		500	{object}	controller.GetRepository.HTTPError
//	@Failure		502	{object}	controller.GetRepository.HTTPError
//	@Router			/subscriptions/{owner}/{repo} [delete]
func (h *Handler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")

	if owner == "" || repo == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid owner or repo")
		return
	}

	if err := h.SubscribeUseCase.Unsubscribe(owner, repo); err != nil {
		h.log.Debug("http: error", "error", err)
		switch apperror.CodeOf(err) {
		case apperror.CodeNotFound:
			writeJSONError(w, http.StatusNotFound, "not found")
		case apperror.CodeInvalidArgument:
			writeJSONError(w, http.StatusBadRequest, "invalid owner or repo")
		case apperror.CodeUnavailable:
			writeJSONError(w, http.StatusBadGateway, "unavailable github service")
		case apperror.CodeDublicate:
			writeJSONError(w, http.StatusConflict, "already subscribed")
		default:
			fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
			writeJSONError(w, http.StatusInternalServerError, "internal error")
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetSubscriptions godoc
//
//	@Summary		Get subscriptions
//	@Description	Get all repository subscriptions
//	@ID				get-subscriptions
//	@Tags			subscriptions
//	@Produce		json
//	@Success		200	{array}		domain.Subscription
//	@Failure		500	{object}	controller.GetRepository.HTTPError
//	@Router			/subscriptions [get]
func (h *Handler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	subs, err := h.SubscribeUseCase.GetSubscriptions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(subs)
}

// GetSubscribedRepositories godoc
//
//	@Summary		Get subscribed repository information
//	@Description	Get GitHub information for every subscribed repository
//	@ID				get-subscribed-repository-info
//	@Tags			subscriptions
//	@Produce		json
//	@Success		200	{array}		domain.Repository
//	@Failure		500	{object}	controller.GetRepository.HTTPError
//	@Failure		502	{object}	controller.GetRepository.HTTPError
//	@Router			/subscriptions/info [get]
func (h *Handler) GetSubscribedRepositories(w http.ResponseWriter, r *http.Request) {
	repos, err := h.RepositoryUseCase.GetSubscribedRepository()
	if err != nil {
		switch apperror.CodeOf(err) {
		case apperror.CodeUnavailable:
			writeJSONError(w, http.StatusBadGateway, "unavailable github service")
		default:
			fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
			writeJSONError(w, http.StatusInternalServerError, "internal error")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(repos)
}
