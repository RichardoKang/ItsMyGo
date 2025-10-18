package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...any) (string, []any)

var generators = make(map[Type]generator)

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderby
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ") //
}

func _insert(values ...any) (string, []any) {
	tableName := values[0]                             //eg: "User"
	fields := strings.Join(values[1].([]string), ", ") //values[1]内容: []string{"Name", "Age", "XXX"}
	return fmt.Sprintf("INSERT INTO %s (%s)", tableName, fields), []any{}
}

func _values(values ...any) (string, []any) {
	var bindStr string //eg: "?, ?, ?"

	var sql strings.Builder
	var vars []any
	sql.WriteString("VALUES ") //eg: "VALUES (?, ?, ?), (?, ?, ?)"

	for i, value := range values {
		v := value.([]any) //把每一行数据都转换成[]any类型
		if bindStr == "" { //如果bindStr还没有生成，就生成一个
			bindStr = genBindVars(len(v)) //比如v是[]any{"Tom", 18}，生成"?, ?"
		}
		sql.WriteString(fmt.Sprintf("(%s)", bindStr)) //把"?,?"变成"(?, ?)"，然后拼接到sql中
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...) //eg: vars = []any{"Tom", 18, "Jack", 20}
	}
	return sql.String(), vars //eg: "VALUES (?, ?), (?, ?)", []any{"Tom", 18, "Jack", 20}
}

func _select(values ...any) (string, []any) {
	tablename := values[0] //eg: "User"
	fields := strings.Join(values[1].([]string), ", ")
	return fmt.Sprintf("SELECT %s FROM %s", fields, tablename), []any{}
}

func _limit(values ...any) (string, []any) {
	return "LIMIT ?", values
}

func _where(values ...any) (string, []any) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _orderby(values ...any) (string, []any) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []any{}
}

func _update(values ...any) (string, []any) {
	tableName := values[0]

	m := values[1].(map[string]any)

	var (
		keys []string
		vars []any
	)

	for key, value := range m {
		keys = append(keys, key+"=?")
		vars = append(vars, value)
	}

	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

func _delete(values ...any) (string, []any) {
	tableName := values[0]

	return fmt.Sprintf("DELETE FROM %s", tableName), []any{}
}

func _count(values ...any) (string, []any) {
	tableName := values[0]
	return _select(tableName, []string{"count(*)"})
}
