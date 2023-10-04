package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func Table(header []string, data [][]string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("  ") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	return table
}

type Render interface {
	WriteLine(data interface{})
	Render()
}

type render struct {
	table *tablewriter.Table
	deal  func(interface{}) []string
	line  int
}

func newRender(header []string, deal func(interface{}) []string) *render {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("  ") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.SetHeader(header)
	return &render{
		table: table,
		deal:  deal,
	}
}

func (e *render) WriteLine(data interface{}) {
	e.line += 1
	e.table.Append(e.deal(data))
}

func (e *render) Render() {
	if e.line == 0 {
		fmt.Println("没有数据")
	} else {
		e.table.Render()
	}
}

func DealTable(header []string, data interface{}, deal func(interface{}) []string) Render {
	r := newRender(header, deal)
	reflectionRow := reflect.ValueOf(data)
	if reflectionRow.Kind() != reflect.Slice && reflectionRow.Kind() != reflect.Array {
		return r
	}
	for i := 0; i < reflectionRow.Len(); i++ {
		r.WriteLine(reflectionRow.Index(i).Interface())
	}
	return r
}

func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// Abs 转绝对路径,处理了"~"
func Abs(p string) (string, error) {
	if strings.HasPrefix(p, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Abs(filepath.Join(home, p[1:]))
	}
	return filepath.Abs(p)
}
