package tuna

// A Row maps a column name to a raw value.
type Row map[string]string

// Set set the field of a Row and then returns the Row.
func (r Row) Set(k string, v string) Row {
	r[k] = v
	return r
}

// An ErrRow is a Row that has an accompanying error.
type ErrRow struct {
	Row
	err error
}
