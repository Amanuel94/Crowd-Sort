package io

import (
	"fmt"
	"os"
	"os/exec"
	"text/tabwriter"

	"github.com/Amanuel94/crowdsort/shared"
	"github.com/TreyBastian/colourize"
)

// for printing the leaderboard as a table

func printTable[T any](header []string, data []shared.Wire[T], p shared.Connector[T]) {

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)

	fmt.Fprintln(writer, formatHeader(header))
	fmt.Fprintln(writer, formatSeparator(len(header)))

	for _, row := range data {
		fmt.Fprintln(writer, formatRow(row, p))
	}

	writer.Flush()
}

func formatHeader(header []string) string {
	return fmt.Sprintf("%s\t%v\t", header[0], header[1])
}

func formatRow[T any](item shared.Wire[T], p shared.Connector[T]) string {
	color := colourize.White
	if item.GetIndex() == p.F || item.GetIndex() == p.S {
		color = colourize.Green
	}
	return fmt.Sprintf("%s\t%v\t", colourize.Colourize(item.GetIndex(), color), colourize.Colourize(item.GetValue(), color))
	// return fmt.Sprintf("%s\t%v\t", item.GetIndex(), item.GetValue())
}

func printProgressBar(current, total int) {
	width := 50
	progress := float64(current) / float64(total)
	barWidth := int(progress * float64(width))

	bar := "[" + string(repeat('#', barWidth)) + string(repeat(' ', width-barWidth)) + "]"
	percentage := int(progress * 100)

	fmt.Printf("\r%s %3d%%", bar, percentage)
}

func repeat(char rune, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += string(char)
	}
	return result
}

func formatSeparator(columns int) string {
	separator := ""
	for i := 0; i < columns; i++ {
		separator += "--------\t"
	}
	return separator
}

func clearTable() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
