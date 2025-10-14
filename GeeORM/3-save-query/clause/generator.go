package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators = make(map[Type]generator)

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderby
}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ") //
}

func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0]                             //eg: "User"
	fields := strings.Join(values[1].([]string), ", ") //values[1]内容: []string{"Name", "Age", "XXX"}
	return fmt.Sprintf("INSERT INTO %s (%s)", tableName, fields), []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	var bindStr string //eg: "?, ?, ?"

	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ") //eg: "VALUES (?, ?, ?), (?, ?, ?)"

	for i, value := range values {
		v := value.([]interface{}) //把每一行数据都转换成[]interface{}类型
		if bindStr == "" {         //如果bindStr还没有生成，就生成一个
			bindStr = genBindVars(len(v)) //比如v是[]interface{}{"Tom", 18}，生成"?, ?"
		}
		sql.WriteString(fmt.Sprintf("(%s)", bindStr)) //把"?,?"变成"(?, ?)"，然后拼接到sql中
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...) //eg: vars = []interface{}{"Tom", 18, "Jack", 20}
	}
	return sql.String(), vars //eg: "VALUES (?, ?), (?, ?)", []interface{}{"Tom", 18, "Jack", 20}
}

func _select(values ...interface{}) (string, []interface{}) {
	tablename := values[0] //eg: "User"
	fields := strings.Join(values[1].([]string), ", ")
	return fmt.Sprintf("SELECT %s FROM %s", fields, tablename), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _orderby(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}
