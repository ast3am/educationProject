package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/ast3am/educationProject/internal/models"
	"github.com/ast3am/educationProject/pkg/logging"
	"strconv"
)

type repository struct {
	storage map[string]*models.UserModel
	id      int
	logger  *logging.Logger
}

func NewRepository(ctx context.Context, rep map[string]*models.UserModel, logger *logging.Logger) *repository {
	return &repository{
		storage: rep,
		logger:  logger,
	}
}

func (r *repository) Create(ctx context.Context, user *models.UserModel) error {
	r.storage[user.ID] = user
	r.logger.Debug().Msg("method Create finished")
	return nil
}

func (r *repository) MakeFriends(ctx context.Context, id, id2 string) (string, error) {
	var err error
	// проверка на существование пользователей
	_, ok := r.storage[id]
	_, ok2 := r.storage[id2]
	switch {
	case !ok && !ok2:
		{
			err = errors.New("Пользователи " + id + " " + id2 + " не найдены\n")
		}
	case !ok:
		{
			err = errors.New("Пользователь " + id + " не найден\n")
		}
	case !ok2:
		{
			err = errors.New("Пользователь " + id2 + " не найден\n")
		}
	}

	if err != nil {
		return "", err
	}

	// проверка, не являются ли друзьями
	for _, v := range r.storage[id].Friends {
		if v == r.storage[id2] {
			err = errors.New("Пользователи " + id + " " + id2 + " уже друзья\n")
		}
	}

	if err != nil {
		return "", err
	}

	// добавление в друзья
	r.storage[id].Friends = append(r.storage[id].Friends, r.storage[id2])
	r.storage[id2].Friends = append(r.storage[id2].Friends, r.storage[id])
	r.logger.Debug().Msgf("method MakeFriends finished + %v", r.storage[id])
	return fmt.Sprint(r.storage[id].Name, " и ", r.storage[id2].Name, " теперь друзья"), nil
}
func (r *repository) Delete(ctx context.Context, id string) (string, error) {
	//проверка на существование
	_, ok := r.storage[id]
	if !ok {
		err := errors.New("Пользователь " + id + " не найден\n")
		return "", err
	}

	//удаление из друзей
	friendArr := make([]string, 0)
	for _, friends := range r.storage[id].Friends {
		friendArr = append(friendArr, friends.ID)
	}
	for _, some := range friendArr {
		for i, v := range r.storage[some].Friends {
			if r.storage[id] == v {
				r.storage[some].Friends = append(r.storage[some].Friends[0:i], r.storage[some].Friends[i+1:]...)
			}
		}
	}
	name := r.storage[id].Name

	//удаление из хранилища
	delete(r.storage, id)
	r.logger.Debug().Msg("method Delete finished")
	return fmt.Sprint("пользователь ", name, " удален"), nil
}

func (r *repository) FindFriend(ctx context.Context, id string) (ufriends []*models.UserModel, err error) {
	//проверка на существование
	_, ok := r.storage[id]
	if !ok {
		err := errors.New("Пользователь " + id + " не найден\n")
		return nil, err
	}
	// передача списка друзей
	ufriends = r.storage[id].Friends
	r.logger.Debug().Msg("method FindFriend finished")
	return
}
func (r *repository) UpdateAge(ctx context.Context, id, age string) error {
	//проверка на существование
	_, ok := r.storage[id]
	if !ok {
		err := errors.New("Пользователь " + id + " не найден\n")
		return err
	}
	r.storage[id].Age = age
	r.logger.Debug().Msg("method UpdateAge finished")
	return nil
}

func (r *repository) MakeID() string {
	r.id++
	id := strconv.Itoa(r.id)
	return id
}
