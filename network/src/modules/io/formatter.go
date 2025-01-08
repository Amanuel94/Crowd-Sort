package io

import (
	"fmt"
	"network/shared"
	"os"
	"text/tabwriter"
)

// for printing the leaderboard as a table

func printTable[T any](header []string, data []shared.IndexedItem[T]) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)

	fmt.Fprintln(writer, formatHeader(header))
	fmt.Fprintln(writer, formatSeparator(len(header)))

	// Print rows
	for _, row := range data {
		fmt.Fprintln(writer, formatRow(row))
	}

	// Flush the writer
	writer.Flush()
}

func formatHeader(header []string) string {
	return fmt.Sprintf("%s\t%s\t", header[0], header[1])
}

func formatRow[T any](item shared.IndexedItem[T]) string {
	return fmt.Sprintf("%s\t%s\t", item.GetIndex(), item.GetValue())
}

func formatSeparator(columns int) string {
	separator := ""
	for i := 0; i < columns; i++ {
		separator += "--------\t"
	}
	return separator
}
