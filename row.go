package tuna

// A Row maps a column name to a raw value.
type Row map[string]string

// An ErrRow is a Row that has an accompanying error.
type ErrRow struct {
	Row
	Err error
}
