package todolist_test

import (
	"errors"
	"fmt"
	"github.com/go-bumbu/todo-app/internal/model/todolist"
	"github.com/google/go-cmp/cmp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
)

func TestListTasks(t *testing.T) {

	tmpDir := t.TempDir()
	db, err := gorm.Open(sqlite.Open(filepath.Join(tmpDir, "test.db")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}

	mngr, err := todolist.New(db)
	if err != nil {
		t.Fatal(err)
	}

	const User1 = "u1"
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 20; i++ {
			_ = createTask(t, mngr, "task"+strconv.Itoa(i)+"_"+User1, User1)
		}
	}()

	const User2 = "u2"
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 30; i++ {
			_ = createTask(t, mngr, "task"+strconv.Itoa(i)+"_"+User2, User2)
		}
	}()
	wg.Wait()

	t.Run("head of index", func(t *testing.T) {
		items, err := mngr.List(User1, 2, 0)
		if err != nil {
			t.Fatal(err)
		}
		got := []string{}
		for _, item := range items {
			got = append(got, item.Text)
		}

		want := []string{
			"task1_u1",
			"task2_u1",
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}
	})

	t.Run("different user and page", func(t *testing.T) {
		items, err := mngr.List(User2, 3, 2)
		if err != nil {
			t.Fatal(err)
		}
		got := []string{}
		for _, item := range items {
			got = append(got, item.Text)
		}

		want := []string{
			"task4_u2",
			"task5_u2",
			"task6_u2",
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}
	})
}

const inMemorySqlite = "file::memory:?cache=shared"

func TestCrudTask(t *testing.T) {
	// initialize DB
	db, err := gorm.Open(sqlite.Open(inMemorySqlite), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	mngr, err := todolist.New(db)
	if err != nil {
		t.Fatal(err)
	}
	// create a bunch of tasks for different users
	t1 := createTask(t, mngr, "task1", "u1")
	t2 := createTask(t, mngr, "task2", "u1")
	t3 := createTask(t, mngr, "task1", "u2")

	// verify we can read the tasks
	readTask(t, mngr, t1, "u1", "task1", "")
	// make sure ownership is conserved
	readTask(t, mngr, t1, "u2", "", fmt.Sprintf("task with id: %s and owner u2 not found", t1))
	readTask(t, mngr, t2, "u1", "task2", "")
	readTask(t, mngr, t3, "u2", "task1", "")
	readTask(t, mngr, t3, "u3", "", fmt.Sprintf("task with id: %s and owner u3 not found", t3))

	// Complete a TodoItem
	setDone(t, mngr, t1, "u1", true, "")
	// send complete again
	setDone(t, mngr, t1, "u1", true, "")
	// wrong owner
	setDone(t, mngr, t1, "u2", true, fmt.Sprintf("task with id: %s and owner u2 not found", t1))
	// set to pending
	setDone(t, mngr, t1, "u1", false, "")
	// again
	setDone(t, mngr, t1, "u1", false, "")

	// update the text
	setText(t, mngr, t1, "u1", "task1Updated", "")
	setText(t, mngr, t1, "u2", "", fmt.Sprintf("task with id: %s and owner u2 not found", t1))

	// update more than one
	setMultiple(t, mngr, t1, "u1", todolist.TodoItem{Text: "multiUpdate", Done: true}, "")
	setMultiple(t, mngr, t1, "u2", todolist.TodoItem{Text: "multiUpdate", Done: true}, fmt.Sprintf("task with id: %s and owner u2 not found", t1))

	// delete the TodoItem
	deleteTask(t, mngr, t1, "u2", fmt.Sprintf("task with id: %s and owner u2 not found", t1))
	deleteTask(t, mngr, t1, "u1", "")

}

func createTask(t *testing.T, mngr *todolist.Manager, content, owner string) string {
	task := todolist.TodoItem{
		Text:    content,
		OwnerId: owner,
	}
	id, err := mngr.Create(&task)
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Error("returned id should not be empty")
	}
	return id
}

func readTask(t *testing.T, mngr *todolist.Manager, taskId, owner, want string, wantErr string) {
	// verify we can read the tasks
	task, err := mngr.Get(taskId, owner)
	if err != nil {
		if wantErr != err.Error() {
			t.Errorf("wanted error:\"%s\", but got: \"%s\"", wantErr, err.Error())
		}
	}
	if task.Text != want {
		t.Errorf("expect TodoItem value to be \"%s\", but got: \"%s\"", want, task.Text)
	}
}

func setDone(t *testing.T, mngr *todolist.Manager, taskId, owner string, val bool, wantErr string) {
	var err error
	err = mngr.Update(taskId, owner, "", &val)
	if err != nil {
		if wantErr != err.Error() {
			t.Errorf("wanted error:\"%s\", but got: \"%s\"", wantErr, err.Error())
		}
		return
	}
	task, err := mngr.Get(taskId, owner)
	if err != nil {
		t.Error(err)
	}
	if task.Done != val {
		t.Errorf("expect TodoItem value to be \"%t\", but got: \"%t\"", val, task.Done)
	}
}

func setText(t *testing.T, mngr *todolist.Manager, taskId, owner, text, wantErr string) {
	err := mngr.Update(taskId, owner, text, nil)
	if err != nil {
		if wantErr != err.Error() {
			t.Errorf("wanted error:\"%s\", but got: \"%s\"", wantErr, err.Error())
		}
		return
	}
	task, err := mngr.Get(taskId, owner)
	if err != nil {
		t.Error(err)
	}
	if task.Text != text {
		t.Errorf("expect TodoItem value to be \"%s\", but got: \"%s\"", text, task.Text)
	}
}

func setMultiple(t *testing.T, mngr *todolist.Manager, taskId, owner string, input todolist.TodoItem, wantErr string) {
	err := mngr.Update(taskId, owner, input.Text, &input.Done)
	if err != nil {
		if wantErr != err.Error() {
			t.Errorf("wanted error:\"%s\", but got: \"%s\"", wantErr, err.Error())
		}
		return
	}
	task, err := mngr.Get(taskId, owner)
	if err != nil {
		t.Error(err)
	}
	if task.Text != input.Text {
		t.Errorf("expect TodoItem value to be \"%s\", but got: \"%s\"", input.Text, task.Text)
	}
	if task.Done != input.Done {
		t.Errorf("expect TodoItem value to be \"%t\", but got: \"%t\"", input.Done, task.Done)
	}
}

func deleteTask(t *testing.T, mngr *todolist.Manager, taskId, owner, wantErr string) {
	err := mngr.Delete(taskId, owner)
	if err != nil {
		if wantErr != err.Error() {
			t.Errorf("wanted error:\"%s\", but got: \"%s\"", wantErr, err.Error())
		}
		return
	}
	_, err = mngr.Get(taskId, owner)
	if err != nil {
		target := &todolist.ItemNotFountErr{}
		if !errors.As(err, &target) {
			t.Errorf("unexpected error: %v", err)
		}

	}

}
