package sq

import "bytes"

type RowValue []interface{}

func (r RowValue) AppendSQLExclude(dialect string, buf *bytes.Buffer, args *[]interface{}, params map[string]int, excludedTableQualifiers []string) error {
	buf.WriteString("(")
	var err error
	for i, value := range r {
		if i > 0 {
			buf.WriteString(", ")
		}
		err = appendSQLValue(buf, args, params, excludedTableQualifiers, value)
		if err != nil {
			return err
		}
	}
	buf.WriteString(")")
	return nil
}

func (r RowValue) GetName() string  { return "" }
func (r RowValue) GetAlias() string { return "" }

func (r RowValue) AppendSQL(dialect string, buf *bytes.Buffer, args *[]interface{}, params map[string]int) error {
	return r.AppendSQLExclude(dialect, buf, args, params, nil)
}

func (r RowValue) In(v interface{}) CustomPredicate {
	if v, ok := v.(RowValue); ok {
		return Predicatef("? IN ?", r, v)
	}
	return Predicatef("? IN (?)", r, v)
}

type RowValues []RowValue

func (rs RowValues) AppendSQL(dialect string, buf *bytes.Buffer, args *[]interface{}, params map[string]int) error {
	var err error
	for i, rowvalue := range rs {
		if i > 0 {
			buf.WriteString(", ")
		}
		err = rowvalue.AppendSQL(dialect, buf, args, params)
		if err != nil {
			return err
		}
	}
	return nil
}
