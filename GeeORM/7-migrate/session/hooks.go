package session

import (
	"geeorm/log"
	"reflect"
)

const (
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
)

func (s *Session) CallMethod(method string, value any) {
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method) // MethodByName 获取该对象的方法

	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method) // 获取传入对象的方法
	}

	param := []reflect.Value{reflect.ValueOf(s)} // 构造参数列表，比如传入当前的 Session 对象

	if fm.IsValid() {
		if v := fm.Call(param); len(v) > 0 { // 调用方法并获取返回值
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}
