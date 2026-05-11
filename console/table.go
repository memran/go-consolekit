package console

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type TableBuilder struct {
	t table.Writer
}

func NewTableBuilder() *TableBuilder {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	return &TableBuilder{t: t}
}

func (tb *TableBuilder) Headers(headers ...string) *TableBuilder {
	row := make(table.Row, len(headers))
	for i, h := range headers {
		row[i] = h
	}
	tb.t.AppendHeader(row)
	return tb
}

func (tb *TableBuilder) Row(values ...string) *TableBuilder {
	row := make(table.Row, len(values))
	for i, v := range values {
		row[i] = v
	}
	tb.t.AppendRow(row)
	return tb
}

func (tb *TableBuilder) Render() {
	tb.t.Render()
}
