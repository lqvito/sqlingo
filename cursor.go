package sqlingo

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
)

type Cursor interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
}

type cursor struct {
	rows *sql.Rows
}

func (c cursor) Next() bool {
	return c.rows.Next()
}

func preparePointers(val reflect.Value, scans *[]interface{}) error {
	kind := val.Kind()
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		*scans = append(*scans, val.Addr().Interface())
	case reflect.Slice:
	case reflect.Struct:
		for j := 0; j < val.NumField(); j++ {
			field := val.Field(j)
			if field.Kind() == reflect.Interface {
				continue
			}
			*scans = append(*scans, field.Addr().Interface())
		}
	case reflect.Ptr:
		toType := val.Type().Elem()
		switch toType.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.String:
			*scans = append(*scans, val.Addr().Interface())
		default:
			to := reflect.New(toType).Elem()
			val.Set(to.Addr())
			err := preparePointers(to, scans)
			if err != nil {
				return nil
			}
		}
	default:
		return fmt.Errorf("unknown type %s", kind.String())
	}
	return nil
}

func parseBool(s []byte) (bool, error) {
	if len(s) == 1 {
		if s[0] == 0 {
			return false, nil
		} else if s[0] == 1 {
			return true, nil
		}
	}
	return strconv.ParseBool(string(s))
}

func (c cursor) Scan(dest ...interface{}) error {
	if len(dest) == 0 {
		// dry run
		return nil
	}

	var scans []interface{}
	for i, item := range dest {
		if reflect.ValueOf(item).Kind() != reflect.Ptr {
			return fmt.Errorf("argument %d is not pointer", i)
		}

		val := reflect.Indirect(reflect.ValueOf(item))

		err := preparePointers(val, &scans)
		if err != nil {
			return err
		}
	}

	pbs := make(map[int]*bool)
	ppbs := make(map[int]**bool)

	for i, scan := range scans {
		if pb, ok := scan.(*bool); ok {
			var s []uint8
			scans[i] = &s
			pbs[i] = pb
		} else if ppb, ok := scan.(**bool); ok {
			var s *[]uint8
			scans[i] = &s
			ppbs[i] = ppb
		}
	}

	err := c.rows.Scan(scans...)
	if err != nil {
		return err
	}

	for i, pb := range pbs {
		if *(scans[i].(*[]byte)) == nil {
			return fmt.Errorf("field %d is null", i)
		}
		b, err := parseBool(*(scans[i].(*[]byte)))
		if err != nil {
			return err
		}
		*pb = b
	}
	for i, ppb := range ppbs {
		if *(scans[i].(**[]uint8)) == nil {
			*ppb = nil
		} else {
			b, err := parseBool(**(scans[i].(**[]byte)))
			if err != nil {
				return err
			}
			*ppb = &b
		}
	}

	return err
}

func (c cursor) Close() error {
	return c.rows.Close()
}
