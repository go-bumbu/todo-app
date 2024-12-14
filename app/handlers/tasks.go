package handlrs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-bumbu/todo-app/internal/model/todolist"
	"github.com/go-bumbu/userauth/handlers/sessionauth"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var _ = spew.Dump // prevent IDE from removing dependency

type TodoListHandler struct {
	TaskManager *todolist.Manager
}

type localTaskList struct {
	Count int
	Tasks []localTaskOutput
}

const limitParam = "limit"
const pageParam = "page"

func (h *TodoListHandler) List() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uData, err := sessionauth.CtxGetUserData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		limitStr := r.URL.Query().Get(limitParam)
		limit := 0
		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				http.Error(w, fmt.Sprintf("unable to convert limit value to number"), http.StatusBadRequest)
				return
			}
		}

		pageStr := r.URL.Query().Get(pageParam)
		page := 0
		if pageStr != "" {
			page, err = strconv.Atoi(pageStr)
			if err != nil {
				http.Error(w, fmt.Sprintf("unable to convert page value to number"), http.StatusBadRequest)
				return
			}
		}

		items, err := h.TaskManager.List(uData.UserId, limit, page)
		if err != nil {
			t := &todolist.ItemNotFountErr{}
			if errors.As(err, &t) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("unable to get task: %s", err.Error()), http.StatusInternalServerError)
			}
			return
		}

		taskItems := make([]localTaskOutput, len(items))
		for i := 0; i < len(items); i++ {

			taskItems[i] = localTaskOutput{
				Id:   items[i].ID,
				Text: items[i].Text,
				Done: items[i].Done,
			}
		}

		output := localTaskList{
			Count: len(taskItems),
			Tasks: taskItems,
		}

		respJson, err := json.Marshal(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respJson)

	})
}

type localTaskInput struct {
	Text string `json:"text"`
	Done *bool
}
type localTaskOutput struct {
	Id   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func (h *TodoListHandler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uData, err := sessionauth.CtxGetUserData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Body == nil {
			http.Error(w, fmt.Sprintf("request had empty body"), http.StatusBadRequest)
			return
		}
		payload := localTaskInput{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to decode json: %s", err.Error()), http.StatusBadRequest)
			return
		}

		if payload.Text == "" {
			http.Error(w, "text cannot be empty req task payload", http.StatusBadRequest)
			return
		}

		if payload.Done == nil {
			f := false
			payload.Done = &f
		}

		t := todolist.TodoItem{
			Text:    payload.Text,
			Done:    *payload.Done,
			OwnerId: uData.UserId,
		}
		Id, err := h.TaskManager.Create(&t)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to store task in DB: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		output := localTaskOutput{
			Id:   Id,
			Text: payload.Text,
			Done: *payload.Done,
		}
		respJson, err := json.Marshal(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respJson)
	})
}

func (h *TodoListHandler) Read() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		taskId, hErr := getTaskId(r)
		if hErr != nil {
			http.Error(w, hErr.Error, hErr.Code)
			return
		}

		uData, err := sessionauth.CtxGetUserData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		Task, err := h.TaskManager.Get(taskId, uData.UserId)
		if err != nil {
			t := &todolist.ItemNotFountErr{}
			if errors.As(err, &t) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("unable to get task: %s", err.Error()), http.StatusInternalServerError)
			}
			return
		}
		output := localTaskOutput{
			Id:   Task.ID,
			Text: Task.Text,
			Done: Task.Done,
		}
		respJson, err := json.Marshal(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respJson)
	})
}

func (h *TodoListHandler) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		taskId, hErr := getTaskId(r)
		if hErr != nil {
			http.Error(w, hErr.Error, hErr.Code)
			return
		}

		uData, err := sessionauth.CtxGetUserData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Body == nil {
			http.Error(w, fmt.Sprintf("request had empty body"), http.StatusBadRequest)
			return
		}
		payload := localTaskInput{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to decode json: %s", err.Error()), http.StatusBadRequest)
			return
		}

		taskText := ""
		if payload.Text == "" {
			taskText = payload.Text
		}

		var taskDone *bool
		if payload.Done != nil {
			taskDone = payload.Done
		}

		err = h.TaskManager.Update(taskId, uData.UserId, taskText, taskDone)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to store task in DB: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	})
}

func (h *TodoListHandler) Delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		taskId, hErr := getTaskId(r)
		if hErr != nil {
			http.Error(w, hErr.Error, hErr.Code)
			return
		}

		uData, err := sessionauth.CtxGetUserData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.TaskManager.Delete(taskId, uData.UserId)
		if err != nil {
			t := &todolist.ItemNotFountErr{}
			if errors.As(err, &t) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("unable to get task: %s", err.Error()), http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusAccepted)
	})
}

type httpErr struct {
	Error string
	Code  int
}

func getTaskId(r *http.Request) (string, *httpErr) {
	vars := mux.Vars(r)
	taskId, ok := vars["ID"]
	if !ok {
		return "", &httpErr{
			Error: "could not extract id to read from request context",
			Code:  http.StatusInternalServerError,
		}
	}
	if taskId == "" {
		return "", &httpErr{
			Error: "no task id provided",
			Code:  http.StatusBadRequest,
		}
	}
	_, err := uuid.Parse(taskId)
	if err != nil {
		return "", &httpErr{
			Error: "task id is not a UUID",
			Code:  http.StatusBadRequest,
		}
	}
	return taskId, nil
}
