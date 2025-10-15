package session

import (
	"geeorm/clause"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)

	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames) // 设置 INSERT 子句
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...) // 设置 VALUES 子句
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected() // 返回插入的记录数, nil
}

func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))                // 获取切片的反射值
	destType := destSlice.Type().Elem()                                   // 获取切片元素的类型
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable() // 获取切片元素类型的表信息

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)                              // 设置 SELECT 子句
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT) // 构建 SQL 语句
	rows, err := s.Raw(sql, vars...).QueryRows()                                           // 执行查询
	if err != nil {
		return err
	}

	// 遍历结果集
	for rows.Next() {
		dest := reflect.New(destType).Elem() // 创建切片元素类型的实例
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface()) // 获取字段地址
		}
		if err := rows.Scan(values...); err != nil { // 扫描行数据到字段
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest)) // 将实例添加到切片中
	}
	return rows.Close()
}
