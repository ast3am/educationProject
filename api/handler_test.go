package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ast3am/educationProject/api/mocks"
	"github.com/ast3am/educationProject/internal/models"
	"github.com/ast3am/educationProject/pkg/logging"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Create(t *testing.T) {

	testModel := models.UserModel{
		"1",
		"Helen",
		"18",
		[]*models.UserModel{},
	}

	testTable := []struct {
		name                string
		inputBody           string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			"positive",
			`{"name":"Helen","age":"18","friends":[]}`,
			http.StatusCreated,
			"New user created with id:1",
		},
		{
			"negative",
			`{"name":"Helen","age":18,"friends":[]}`,
			http.StatusBadRequest,
			"Unmarshal error \njson: cannot unmarshal number into Go struct field UserModel.age of type string",
		},
	}

	ctx := context.Background()
	log := logging.GetLogger()
	repository := mocks.NewRepository(t)

	handler := NewHandler(repository, log)

	for _, test := range testTable {
		if test.name == "positive" {
			repository.
				On("MakeID").Return("1").
				On("Create", ctx, &testModel).Return(nil)
		}
		var jsonStr = []byte(test.inputBody)
		req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonStr))
		if err != nil {
			fmt.Errorf("err %+v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code != test.expectedStatusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				w.Code, test.expectedStatusCode)
		}

		if w.Body.String() != test.expectedRequestBody {
			t.Errorf("handler returned unexpected body: got %v want %v",
				w.Body.String(), test.expectedRequestBody)
		}
	}
}
func TestHandler_MakeFriends(t *testing.T) {

	testTable := []struct {
		name                string
		inputBody           string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			"positive",
			`{"source_id":"1","target_id":"2"}`,
			http.StatusOK,
			"Пользователи 1 и 2 теперь друзья",
		},
		{
			"negative1",
			`{"source_id":"1","target_id":""}`,
			http.StatusBadRequest,
			"Unmarshal error: some ID is nil",
		},
		{
			"negative2",
			`{"source_id":"1","target_id":2}`,
			http.StatusBadRequest,
			"Unmarshal error \njson: cannot unmarshal number into Go struct field MakeFriendRequest.target_id of type string",
		},
	}

	ctx := context.Background()
	log := logging.GetLogger()
	repository := mocks.NewRepository(t)

	handler := NewHandler(repository, log)

	for _, test := range testTable {
		if test.name == "positive" {
			repository.
				On("MakeFriends", ctx, "1", "2").Return(test.expectedRequestBody, nil)
		}
		var jsonStr = []byte(test.inputBody)
		req, err := http.NewRequest("POST", "/make_friends", bytes.NewBuffer(jsonStr))
		if err != nil {
			fmt.Errorf("err %+v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.MakeFriends(w, req)

		if w.Code != test.expectedStatusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				w.Code, test.expectedStatusCode)
		}

		if w.Body.String() != test.expectedRequestBody {
			t.Errorf("handler returned unexpected body: got %v want %v",
				w.Body.String(), test.expectedRequestBody)
		}
	}
}
func TestHandler_Delete(t *testing.T) {
	testTable := []struct {
		name                string
		inputBody           string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			"positive",
			`{"target_id":"2"}`,
			http.StatusOK,
			"Пользователь 2 удален",
		},
		{
			"negative",
			`{"target_id":2}`,
			http.StatusBadRequest,
			"Unmarshal error \njson: cannot unmarshal number into Go struct field GetID.target_id of type string",
		},
	}

	ctx := context.Background()
	log := logging.GetLogger()
	repository := mocks.NewRepository(t)

	handler := NewHandler(repository, log)

	for _, test := range testTable {
		if test.name == "positive" {
			repository.
				On("Delete", ctx, "2").Return(test.expectedRequestBody, nil)
		}
		var jsonStr = []byte(test.inputBody)
		req, err := http.NewRequest("DELETE", "/user", bytes.NewBuffer(jsonStr))
		if err != nil {
			fmt.Errorf("err %+v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h := http.HandlerFunc(handler.Delete)
		h.ServeHTTP(w, req)

		if w.Code != test.expectedStatusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				w.Code, test.expectedStatusCode)
		}

		if w.Body.String() != test.expectedRequestBody {
			t.Errorf("handler returned unexpected body: got %v want %v",
				w.Body.String(), test.expectedRequestBody)
		}
	}
}
