package ql

import (
	"fmt"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	var sql = "SELECT MEAN(flux), b FROM table"

	var p = NewParser(strings.NewReader(sql))

	fmt.Println(p.Parse())

}
