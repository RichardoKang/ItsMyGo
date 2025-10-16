package session

import (
	"testing"
)

// go::generate go run ../schema/parse.go -struct User
type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	s := NewSession().Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}

func TestSession_Model(t *testing.T) {
	s := NewSession().Model(&User{})
	table := s.RefTable()
	s.Model(&Session{})
	if table.Name != "User" || s.RefTable().Name != "Session" {
		t.Fatal("Failed to change model")
	}
}

func TestSession_Table(t *testing.T) {
	s := NewSession().Model(&User{})
	table := s.RefTable()

	//t.Log("Table Model:", reflect.TypeOf(table.Model).Elem().Name())
	//t.Log("Table Name:", table.Name)
	//t.Log("Table Fields:", table.Fields)

	for _, fieldName := range table.Fields {
		t.Log("Field Name:", fieldName.Name, ";Field Type:", fieldName.Type, ";Field Tag:", fieldName.Tag)
	}
}
