package utils

import (
	"os"
	"reflect"

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

func DealTable(header []string, data interface{}, deal func(interface{}) []string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("  ") // pad with tabs
	table.SetNoWhiteSpace(true)
	reflectionRow := reflect.ValueOf(data)
	if reflectionRow.Kind() != reflect.Slice || reflectionRow.Kind() != reflect.Array {
		return nil
	}
	for i := 0; i < reflectionRow.Len(); i++ {
		table.Append(deal(reflectionRow.Index(i).Interface()))
	}
	return table
}
