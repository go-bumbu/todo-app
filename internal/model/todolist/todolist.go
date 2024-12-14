package todolist

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Manager struct {
	db *gorm.DB
}

func New(db *gorm.DB) (*Manager, error) {
	// Migrate the schema
	err := db.AutoMigrate(&TodoItem{})
	if err != nil {
		return nil, err
	}

	m := Manager{
		db: db,
	}
	return &m, nil
}

type TodoItem struct {
	ID      string `gorm:"primaryKey,index"`
	OwnerId string `gorm:"index"`
	Text    string
	Done    bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (user *TodoItem) BeforeCreate(db *gorm.DB) (err error) {
	// UUID version 4
	user.ID = uuid.NewString()
	return
}

type ItemNotFountErr struct {
	id    string
	owner string
}

func (m *ItemNotFountErr) Error() string {
	return fmt.Sprintf("task with id: %s and owner %s not found", m.id, m.owner)
}

func (m Manager) List(owner string, size, page int) ([]TodoItem, error) {
	if size <= 0 {
		size = 20
	}
	if size >= 50 {
		size = 50
	}

	offset := size * (page - 1)
	if offset <= 0 {
		offset = 0
	}
	tasks := make([]TodoItem, size)
	result := m.db.Where("owner_id = ?", owner).Model(&TodoItem{}).Offset(offset).Limit(size).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

func (m Manager) Create(task *TodoItem) (string, error) {

	result := m.db.Create(task)
	if result.Error != nil {
		return "", result.Error
	}
	return task.ID, nil
}

func (m Manager) Get(id, owner string) (TodoItem, error) {
	t := TodoItem{}
	result := m.db.First(&t, "ID = ? AND owner_id = ?", id, owner)
	if result.RowsAffected == 0 {
		return t, &ItemNotFountErr{id: id, owner: owner}
	}
	return t, nil
}

func (m Manager) Update(id, owner, text string, done *bool) error {

	fieldMap := map[string]any{}
	if text != "" {
		fieldMap["text"] = text
	}
	if done != nil {
		fieldMap["done"] = *done
	}

	t := TodoItem{}
	result := m.db.Model(&t).
		Where("ID = ? AND owner_id = ?", id, owner).
		Updates(fieldMap)

	if result.RowsAffected == 0 {
		return &ItemNotFountErr{id: id, owner: owner}
	}
	return nil
}

func (m Manager) Delete(id, owner string) error {

	t := TodoItem{}
	result := m.db.Where("ID = ? AND owner_id = ?", id, owner).Delete(&t)

	if result.RowsAffected == 0 {
		return &ItemNotFountErr{id: id, owner: owner}
	}
	return nil
}

// TODO, hard delete
