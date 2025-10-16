package clause

import "strings"

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]any
}
type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	WHERE
	LIMIT
	ORDERBY
)

func (c *Clause) Set(name Type, vars ...any) {
	if c.sql == nil {
		c.sql = make(map[Type]string) // corrected type from generator to string
		c.sqlVars = make(map[Type][]any)
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

func (c *Clause) Build(orders ...Type) (string, []any) {
	var sqls []string
	var vars []any
	for _, order := range orders { // orders是一个Type切片，表示SQL语句的各个部分按什么顺序拼接
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)                 //把各个部分的SQL语句按顺序添加到sqls切片中
			vars = append(vars, c.sqlVars[order]...) //把各个部分的SQL变量按顺序添加到vars切片中
		}
	}
	return strings.Join(sqls, " "), vars //eg: "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?", []any{"Tom", 3}
}
