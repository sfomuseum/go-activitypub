package database

import (
	"fmt"
	"strings"
)

func deriveSQLQueryConditions(database_q *Query) (string, []any) {

	where := make([]string, 0)
	args := make([]any, 0)

	if database_q == nil || database_q.Where == nil {
		str_where := "1=1"
		return str_where, args
	}

	for _, c := range database_q.Where.Conditions {

		where = append(where, fmt.Sprintf("%s %s ?", c.Field, c.Operator))
		args = append(args, c.Value)
	}

	// Something something something order by, limit stuff...

	str_where := strings.Join(where, fmt.Sprintf(" %s ", database_q.Where.Relation))
	return str_where, args
}
