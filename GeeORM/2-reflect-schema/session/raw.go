package session

import (
	"database/sql"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/schema"
	"strings"
)

type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	refTable *schema.Schema
	sql      strings.Builder
	sqlVars  []interface{}
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}
func (s *Session) Clear() {
	s.sql.Reset()   // 重置 sql.Builder
	s.sqlVars = nil // 重置 sqlVars 切片
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)                   // 拼接 SQL 语句
	s.sql.WriteString(" ")                   // 添加空格分隔
	s.sqlVars = append(s.sqlVars, values...) // 添加占位符对应的值
	return s
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear() // 无论函数如何退出都要清理掉 SQL 语句和变量
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...) // 返回 *sql.Row
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)

	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)

	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
