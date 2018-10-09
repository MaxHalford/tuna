package tuna

// A Row maps a column name to a raw value.
type Row map[string]string

// Get gets a value and returns an error if the key doesn't exist.
func (r Row) Get(k string) (string, error) {
	v, ok := r[k]
	if ok {
		return v, nil
	}
	return v, ErrUnknownField{k}
}

// Set sets the field of a Row and then returns the Row.
func (r Row) Set(k string, v string) Row {
	r[k] = v
	return r
}

// An ErrRow is a Row that has an accompanying error.
type ErrRow struct {
	Row
	Err error
}
