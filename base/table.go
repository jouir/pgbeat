package base

import "fmt"

// Table represents a table
type Table struct {
	Schema string
	Name   string
}

// NewTable creates a Table object
func NewTable(schema string, name string) *Table {
	return &Table{
		Schema: schema,
		Name:   name,
	}
}

// String converts Table as string
func (t Table) String() string {
	return fmt.Sprintf("%s.%s", QuoteIdent(t.Schema), QuoteIdent(t.Name))
}
