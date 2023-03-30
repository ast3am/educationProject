package main

import (
	"context"
	"fmt"
	"github.com/ast3am/educationProject/api"
	"github.com/ast3am/educationProject/internal/user/db"
	"github.com/ast3am/educationProject/pkg/logging"
	"github.com/ast3am/educationProject/pkg/mongodb"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer func() {
		log.Info().Msg("canceling context")
		cancel()
	}()
	log := logging.GetLogger()
	log.Info().Msg("started")
	router := chi.NewRouter()
	//resultMap := make(map[string]*user.UserModel)
	//repository := db.NewRepository(ctx, resultMap, log)
	mongoDB, err := mongodb.NewClient(ctx, "localhost", "27017", "SomeBase")
	if err != nil {
		fmt.Println("error")
	}
	mongoRepository := db.NewMongoRepository(mongoDB, "1", log)
	handler := api.NewHandler(mongoRepository, log)
	handler.Register(router)
	start(router)

}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, my http service is running"))
}

func start(r chi.Router) {
	r.Get("/", IndexHandler)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
