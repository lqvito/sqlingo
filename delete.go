package sqlingo

import (
	"database/sql"
)

type DeleteWithTable interface {
	Where(conditions ...BooleanExpression) DeleteWithWhere
}

type DeleteWithWhere interface {
	GetSQL() (string, error)
	Execute() (result sql.Result, err error)
}

type deleteStatus struct {
	scope scope
	where BooleanExpression
}

func (d *database) DeleteFrom(table Table) DeleteWithTable {
	return deleteStatus{scope: scope{Database: d, Tables: []Table{table}}}
}

func (s deleteStatus) Where(conditions ...BooleanExpression) DeleteWithWhere {
	s.where = And(conditions...)
	return s
}

func (s deleteStatus) GetSQL() (string, error) {
	whereSql, err := s.where.GetSQL(s.scope)
	if err != nil {
		return "", err
	}
	sqlString := "DELETE FROM " + s.scope.Tables[0].GetSQL(s.scope) + " WHERE " + whereSql

	return sqlString, nil
}

func (s deleteStatus) Execute() (sql.Result, error) {
	sqlString, err := s.GetSQL()
	if err != nil {
		return nil, err
	}
	return s.scope.Database.Execute(sqlString)
}
