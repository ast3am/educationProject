package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ast3am/educationProject/internal/models"
	"github.com/ast3am/educationProject/pkg/logging"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
)

type Repository interface {
	Create(ctx context.Context, user *models.UserModel) error
	MakeFriends(ctx context.Context, sourceId, targetId string) (string, error)
	Delete(ctx context.Context, id string) (string, error)
	FindFriend(ctx context.Context, id string) (ufriends []*models.UserModel, err error)
	UpdateAge(ctx context.Context, id, age string) error
	MakeID() string
}

type handler struct {
	repository Repository
	logger     *logging.Logger
}

func NewHandler(repository Repository, logger *logging.Logger) *handler {
	return &handler{
		repository: repository,
		logger:     logger,
	}
}

func (h *handler) Register(router chi.Router) {
	router.Post("/create", h.Create)
	router.Post("/make_friends", h.MakeFriends)
	router.Delete("/user", h.Delete)
	router.Get("/friends/{id}", h.GetFriends)
	router.Put("/{id}", h.UpdateAge)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.HandlerErrorLog(r, http.StatusInternalServerError, "", err)
		return
	}
	defer r.Body.Close()

	u := models.UserModel{}

	err = json.Unmarshal(content, &u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unmarshal error \n" + err.Error()))
		h.logger.HandlerErrorLog(r, http.StatusBadRequest, "", err)
		return
	}

	u.ID = h.repository.MakeID()
	h.repository.Create(context.TODO(), &u)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("New user created with id:" + u.ID))
	h.logger.HandlerLog(r, http.StatusCreated, "New user created")
}

func (h *handler) MakeFriends(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.HandlerErrorLog(r, http.StatusInternalServerError, "", err)
		return
	}
	defer r.Body.Close()

	// получение ID из тела запроса

	type MakeFriendRequest struct {
		SourceID string `json:"source_id"`
		TargetID string `json:"target_id"`
	}

	mf := MakeFriendRequest{"", ""}

	err = json.Unmarshal(content, &mf)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unmarshal error \n" + err.Error()))
		h.logger.HandlerErrorLog(r, http.StatusBadRequest, "", err)
		return
	}

	// проверка полей

	if mf.SourceID == "" || mf.TargetID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unmarshal error: some ID is nil"))
		h.logger.HandlerErrorLog(r, http.StatusBadRequest, "Unmarshal error: some ID is nil", err)
		return
	}

	// создание друзей

	text, err := h.repository.MakeFriends(context.TODO(), mf.SourceID, mf.TargetID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		h.logger.HandlerErrorLog(r, http.StatusBadRequest, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
	h.logger.HandlerLog(r, http.StatusCreated, "Friends were made")

}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.HandlerErrorLog(r, http.StatusInternalServerError, "", err)
	}
	defer r.Body.Close()

	//получаем ID из запроса

	type GetID struct {
		TargetID string `json:"target_id"`
	}

	id := GetID{""}

	err = json.Unmarshal(content, &id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unmarshal error \n" + err.Error()))
		h.logger.HandlerErrorLog(r, http.StatusBadRequest, "", err)
		return
	}

	// удаляем пользователя

	text, err := h.repository.Delete(context.TODO(), id.TargetID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		h.logger.HandlerErrorLog(r, http.StatusBadRequest, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
	h.logger.HandlerLog(r, http.StatusCreated, "User deleted")
}

func (h *handler) GetFriends(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.New("ID is nil")
		w.Write([]byte("ID is nil"))
		h.logger.HandlerErrorLog(r, http.StatusBadRequest, "", err)
		return
	}

	friends, err := h.repository.FindFriend(context.TODO(), id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		h.logger.HandlerErrorLog(r, http.StatusBadRequest, "", err)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("друзья пользователя с id: " + id + "\n"))
	for _, v := range friends {
		w.Write([]byte(v.Name + "\n"))
	}

	h.logger.HandlerLog(r, http.StatusCreated, "Friends received")
}

func (h *handler) UpdateAge(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.HandlerErrorLog(r, http.StatusInternalServerError, "", err)
	}
	defer r.Body.Close()

	// получение ID из запроса

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.New("ID is nil")
		w.Write([]byte("ID is nil"))
		h.logger.HandlerErrorLog(r, http.StatusBadRequest, "", err)
		return
	}

	// чтение нового возраста из json

	type GetNewAge struct {
		NewAge string `json:"new_age"`
	}

	newAge := GetNewAge{""}

	err = json.Unmarshal(content, &newAge)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unmarshal error \n" + err.Error()))
		h.logger.HandlerErrorLog(r, http.StatusBadRequest, "", err)
		return
	}

	// обновление возраста

	err = h.repository.UpdateAge(context.TODO(), id, newAge.NewAge)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		h.logger.HandlerErrorLog(r, http.StatusBadRequest, "", err)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("пользователь с id: " + id + " обновлен"))
	h.logger.HandlerLog(r, http.StatusCreated, "User updated")
}
