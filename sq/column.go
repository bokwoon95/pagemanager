package sq

import "time"

// NOTE: testing this is better done in the query builder packages

type ColumnMode int

const (
	ColumnModeInsert ColumnMode = 0
	ColumnModeUpdate ColumnMode = 1
)

type Column struct {
	// mode determines if INSERT or UPDATE
	mode ColumnMode
	// INSERT
	rowStart      bool
	rowEnd        bool
	firstField    string
	insertColumns Fields
	rowValues     RowValues
	// UPDATE
	assignments Assignments
}

func NewColumn(mode ColumnMode) *Column {
	var col Column
	col.mode = mode
	return &col
}
func ColumnInsertResult(col *Column) (Fields, RowValues) { return col.insertColumns, col.rowValues }
func ColumnUpdateResult(col *Column) Assignments         { return col.assignments }

func (col *Column) Set(field Field, value interface{}) {
	if field == nil {
		// should I panic with an error here instead?
		return
	}
	// UPDATE mode
	if col.mode == ColumnModeUpdate {
		col.assignments = append(col.assignments, Assign(field, value))
		return
	}
	// INSERT mode
	name := field.GetName()
	if !col.rowStart {
		col.rowStart = true
		col.firstField = name
		col.insertColumns = append(col.insertColumns, field)
		col.rowValues = append(col.rowValues, RowValue{value})
		return
	}
	switch name {
	case col.firstField: // Start a new RowValue
		if !col.rowEnd {
			col.rowEnd = true
		}
		col.rowValues = append(col.rowValues, RowValue{value})
	default: // Append to last RowValue
		if !col.rowEnd {
			col.insertColumns = append(col.insertColumns, field)
		}
		last := len(col.rowValues) - 1
		col.rowValues[last] = append(col.rowValues[last], value)
	}
}

func (col *Column) SetBool(field BooleanField, value bool)      { col.Set(field, value) }
func (col *Column) SetFloat64(field NumberField, value float64) { col.Set(field, value) }
func (col *Column) SetInt(field NumberField, value int)         { col.Set(field, value) }
func (col *Column) SetInt64(field NumberField, value int64)     { col.Set(field, value) }
func (col *Column) SetString(field StringField, value string)   { col.Set(field, value) }
func (col *Column) SetTime(field TimeField, value time.Time)    { col.Set(field, value) }
