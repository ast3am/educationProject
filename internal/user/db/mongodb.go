package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/ast3am/educationProject/internal/models"
	"github.com/ast3am/educationProject/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type db struct {
	collection *mongo.Collection
	id         int
	logger     *logging.Logger
}

func NewMongoRepository(database *mongo.Database, collection string, logger *logging.Logger) *db {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (d *db) Create(ctx context.Context, user *models.UserModel) error {
	_, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return errors.New("error to insert user")
	}
	d.logger.Debug().Msg("User created with id " + user.ID)
	return nil
}

func (d *db) MakeFriends(ctx context.Context, sourceId, targetId string) (string, error) {
	var err error
	var result bson.M
	ok := [2]bool{true, true}
	ids := [2]string{sourceId, targetId}
	// проверка на существование пользователей
	for i, id := range ids {
		err = d.collection.FindOne(ctx, bson.D{{"id", id}}).Decode(&result)
		if err != nil {
			ok[i] = false
		}
	}

	switch {
	case !ok[0] && !ok[1]:
		{
			err = errors.New("Пользователи " + sourceId + " " + targetId + " не найдены\n")
		}
	case !ok[0]:
		{
			err = errors.New("Пользователь " + sourceId + " не найден\n")
		}
	case !ok[1]:
		{
			err = errors.New("Пользователь " + targetId + " не найден\n")
		}
	}
	if err != nil {
		return "", err
	}

	//проверка на друзей
	checkFilter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"id", sourceId}},
				bson.D{{"friends", targetId}},
			}},
	}
	err = d.collection.FindOne(ctx, checkFilter).Decode(&result)
	if err == nil {
		err = errors.New("Пользователи " + sourceId + " " + targetId + " уже друзья\n")
		return "", err
	}

	// обновление друзей в базе
	for i := 0; i <= 1; i++ {
		if i == 1 {
			sourceId, targetId = targetId, sourceId
		}
		updateFilter := bson.D{{"id", sourceId}}
		updateOptions := bson.D{{"$push", bson.D{{"friends", targetId}}}}
		_, err = d.collection.UpdateOne(ctx, updateFilter, updateOptions)
	}
	//перевернем обратно
	sourceId, targetId = targetId, sourceId
	d.logger.Debug().Msgf("method MakeFriends finished with ids %s, %s", sourceId, targetId)

	return fmt.Sprint("пользователи ", sourceId, " и ", targetId, " теперь друзья"), nil
}

func (d *db) Delete(ctx context.Context, id string) (string, error) {
	//поиск по id
	filter := bson.M{"id": id}
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		err = errors.New("failed to execute with filter")
		return "", err
	}

	//проверка на то, что пользователь удален
	if result.DeletedCount == 0 {
		err := errors.New("Пользователь " + id + " не найден")
		return "", err
	}

	// удаление удаленного пользователя из друзей
	updateFilter := bson.D{{"friends", id}}
	updateOptions := bson.D{{"$pull", bson.D{{"friends", id}}}}
	_, err = d.collection.UpdateMany(ctx, updateFilter, updateOptions)

	d.logger.Debug().Msgf("Удален пользователь с id %s", id)
	return fmt.Sprint("пользователь ", id, " удален"), nil
}

func (d *db) FindFriend(ctx context.Context, id string) (ufriends []*models.UserModel, err error) {
	//проверка на существование
	//поиск по id
	check := bson.D{}
	filter := bson.M{"id": id}
	err = d.collection.FindOne(ctx, filter).Decode(&check)
	if err != nil {
		err = errors.New("пользователь с " + id + " не найден")
		return nil, err
	}
	//поиск друзей по id
	var results []bson.D
	friendsFilter := bson.D{{"friends", id}}
	cursor, err := d.collection.Find(ctx, friendsFilter)
	if err = cursor.All(context.TODO(), &results); err != nil {
		d.logger.Err(err).Msg("find results error")
		return nil, err
	}

	//запись друзей в мапу
	for _, result := range results {
		u := models.UserModel{}
		doc, err := bson.Marshal(result)
		if err != nil {
			err = errors.New("can't marshal result")
			return nil, err
		}
		err = bson.Unmarshal(doc, &u)
		if err != nil {
			err = errors.New("can't unmarshal result to struct")
			return nil, err
		}
		ufriends = append(ufriends, &u)
	}
	d.logger.Debug().Msg("method FindFriend finished")
	return ufriends, nil
}

func (d *db) UpdateAge(ctx context.Context, id, age string) error {
	updateFilter := bson.D{{"id", id}}
	updateOptions := bson.D{{"$set", bson.D{{"age", age}}}}
	res, err := d.collection.UpdateOne(ctx, updateFilter, updateOptions)
	if err != nil {
		err = errors.New(fmt.Sprintf("can't update age %v", err))
		return err
	}
	if res.MatchedCount == 0 {
		err = errors.New(fmt.Sprintf("пользователь с id %s не найден", id))
		return err
	}

	d.logger.Debug().Msgf("Обновлен пользователь с id %s", id)
	return nil
}

func (d *db) MakeID() string {
	// генерация ID из базы если идет запрос больше чем от двух сервисов к 1 базе
	var ids string
	var id int

	//сортировка по последнему добавленному id и возвращение последнего элемента из mongo
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"id", -1}}).SetLimit(1)
	cursor, err := d.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		d.logger.Err(err).Msg("Can't get ID from mongo DB")
	}
	//конвертирование результата в мапу
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		d.logger.Err(err).Msg("Can't get ID from mongo DB")
	}

	//если база пустая и не ни одного id
	if results == nil {
		return "1"
	}

	// получение id из мапы, получаем string
	ids = fmt.Sprintf("%v", results[0]["id"])

	// итерация id и перевод обратно в string
	id, _ = strconv.Atoi(ids)
	id++
	ids = strconv.Itoa(id)
	return ids
}
