package sqlingo

import "strings"

type Field interface {
	Expression
}

type NumberField interface {
	NumberExpression
}

type BooleanField interface {
	BooleanExpression
}

type StringField interface {
	StringExpression
}

func newFieldExpression(tableName string, fieldName string) expression {
	shortFieldNameSql := getSQLForName(fieldName)
	fullFieldNameSql := getSQLForName(tableName) + "." + shortFieldNameSql
	return expression{
		builder: func(scope scope) (string, error) {
			if len(scope.Tables) != 1 || scope.lastJoin != nil || scope.Tables[0].GetName() != tableName {
				return fullFieldNameSql, nil
			}
			return shortFieldNameSql, nil
		},
	}
}

func NewNumberField(tableName string, fieldName string) NumberField {
	return newFieldExpression(tableName, fieldName)
}

func NewBooleanField(tableName string, fieldName string) BooleanField {
	return newFieldExpression(tableName, fieldName)
}

func NewStringField(tableName string, fieldName string) StringField {
	return newFieldExpression(tableName, fieldName)
}

type FieldList []Field

func (fields FieldList) GetSQL(scope scope) (string, error) {
	isSingleTable := len(scope.Tables) == 1 && scope.lastJoin == nil
	var sb strings.Builder
	if len(fields) == 0 {
		for i, table := range scope.Tables {
			if i > 0 {
				sb.WriteString(", ")
			}
			if isSingleTable {
				sb.WriteString(table.GetFieldsSQL())
			} else {
				sb.WriteString(table.GetFullFieldsSQL())
			}
		}
	} else {
		fieldsSql, err := commaFields(scope, fields)
		if err != nil {
			return "", err
		}
		sb.WriteString(fieldsSql)
	}
	return sb.String(), nil
}
