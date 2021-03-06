package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	//generator[INSERT] =
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
}

// 获取绑定的参数占位   ?,?,?
func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ",")
}

// insert into user (a, b, c)
func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join((values[1]).([]string), ",")
	return fmt.Sprintf("insert into %s (%v)", tableName, fields), []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	// 拼成=> values ?, ?, ?
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("values ")
	for i, value := range values {
		// 拼凑
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("%v", bindStr))

		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)

	}

	return "", []interface{}{}
}

func _select(values ...interface{}) (string, []interface{}) {
	// select column1, column2 from tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("select %s from %s", fields, tableName), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	// limit desc
	return "limit ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	// where $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("where %s", desc), vars
}
func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("order by %s", values[0]), []interface{}{}
}
