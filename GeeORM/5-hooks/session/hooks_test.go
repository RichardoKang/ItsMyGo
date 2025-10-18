package session

import (
	"geeorm/log"
	"reflect"
	"testing"
)

type Account struct {
	ID       int `geeorm:"PRIMARY KEY"`
	Password string
}

func (account *Account) BeforeInsert(s *Session) error {
	log.Info("before inert", account)
	account.ID += 1000
	return nil
}

func (account *Account) AfterQuery(s *Session) error {
	log.Info("after query", account)
	account.Password = "******"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	s := NewSession().Model(&Account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	_, _ = s.Insert(&Account{1, "123456"}, &Account{2, "qwerty"})

	u := &Account{}

	err := s.First(u)
	if err != nil || u.ID != 1001 || u.Password != "******" {
		t.Fatal("Failed to call hooks after query, got", u)
	}
}

func TestSession_CallMethod_Debug(t *testing.T) {
	s := NewSession().Model(&Account{})

	if err := s.DropTable(); err != nil {
		t.Fatalf("DropTable error: %v", err)
	}
	if err := s.CreateTable(); err != nil {
		t.Fatalf("CreateTable error: %v", err)
	}

	n, err := s.Insert(&Account{1, "123456"}, &Account{2, "qwerty"})
	if err != nil {
		t.Fatalf("Insert error: %v", err)
	}
	t.Logf("Insert returned n=%d", n)

	// 反射检查 hook 是否存在（辅助定位）
	rt := reflect.TypeOf(&Account{})
	if _, ok := rt.MethodByName("BeforeInsert"); !ok {
		t.Fatal("BeforeInsert method not found on *Account")
	}
	if _, ok := rt.MethodByName("AfterQuery"); !ok {
		t.Fatal("AfterQuery method not found on *Account")
	}

	u := &Account{}
	if err := s.First(u); err != nil {
		t.Fatalf("First error: %v", err)
	}
	t.Logf("First returned: %+v", u)

	if u.ID != 1001 || u.Password != "******" {
		t.Fatalf("Failed to call hooks after query, got %+v", u)
	}
}
