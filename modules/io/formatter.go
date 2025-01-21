package io

import (
	"fmt"
	"os"
	"os/exec"
	"text/tabwriter"

	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/shared"
	"github.com/TreyBastian/colourize"
	"github.com/kataras/tablewriter"
	"github.com/lensesio/tableprinter"
)

func printTable[T any](_ []string, data []shared.Wire[T], p shared.PingMessage) {

	defaultWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', tabwriter.Debug)
	fmt.Fprintln(defaultWriter, "\nLive Update")
	newLine(1)
	writer := tabwriter.NewWriter(os.Stdout, 0, 10, 3, ' ', tabwriter.Debug)
	fmt.Fprintf(defaultWriter, "%s", fmt.Sprintf("   %s\t%s\t%s", "Wire", "Value|", "Status"))
	fmt.Fprintln(defaultWriter)
	for _, row := range data {
		fmt.Fprintln(writer, formatRow(row, p))
	}

	defaultWriter.Flush()

	writer.Flush()
	newLine(2)
}

func formatRow[T any](item shared.Wire[T], p shared.PingMessage) string {
	color := colourize.White
	if item.GetIndex() == p.F || item.GetIndex() == p.S {
		color = colourize.Green
	}
	s := colourize.Colourize(item.GetIndex(), color)
	v := colourize.Colourize(item.GetValue(), color)
	statColor := colourize.White
	if item.GetStatus() == shared.COMPLETED {
		statColor = colourize.Green
	} else if item.GetStatus() == shared.PENDING {
		statColor = colourize.Yellow

	}
	stat := colourize.Colourize(item.GetStatus(), statColor)
	return fmt.Sprintf("   %s\t%s\t  \t%s", s, v, stat)
}

func printProgressBar(current, total int) {
	width := 50
	progress := float64(current) / float64(total)
	barWidth := int(progress * float64(width))

	bar := colourize.Colourize("   ["+string(repeat('#', barWidth))+string(repeat(' ', width-barWidth))+"]", colourize.Blue)
	percentage := int(progress * 100)

	fmt.Printf("\r%s %3d%%", bar, percentage)
	newLine(1)
}

func repeat(char rune, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += string(char)
	}
	return result
}

func clearTable() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printUpdate(p shared.PingMessage) {
	msg := colourize.Colourize(fmt.Sprintf("   Comparator %s submitted a task.\n    F: %s\n    S: %s", p.AssignieeId, p.F, p.S), colourize.Green, colourize.Bold)
	fmt.Println(msg)
	newLine(1)
}

type comparatorRow struct {
	Comparator string `header:"Comparator"`
	TaskCount  string `header:"Task Count"`
	Status     string `header:"Status"`
}

func printWorkerStatusTable[T any](workers []*(interfaces.Comparator[T])) {

	printer := tableprinter.New(os.Stdout)
	tble := make([]comparatorRow, 0)
	for _, worker := range workers {
		row := comparatorRow{
			Comparator: shared.AsModule(worker).GetID().(string),
			TaskCount:  fmt.Sprintf("%v", shared.AsModule(worker).TaskCount()),
			Status:     shared.AsModule(worker).GetStatus(),
		}
		tble = append(tble, row)
	}

	printer.BorderTop, printer.BorderLeft, printer.BorderRight = true, true, true
	printer.CenterSeparator = " "
	printer.ColumnSeparator = " "
	printer.RowSeparator = "─"
	printer.HeaderFgColor = tablewriter.FgGreenColor // set header foreground color for all headers.
	printer.Print(tble)

	newLine(2)

}
