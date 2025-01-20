package io

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/shared"
	"github.com/TreyBastian/colourize"
)

func printTable[T any](_ []string, data []shared.Wire[T], p shared.PingMessage) {

	defaultWriter := tabwriter.NewWriter(os.Stdout, 0, 2, 10, ' ', tabwriter.AlignRight)
	fmt.Fprintln(defaultWriter, "\n\tLive Update")
	fmt.Fprintln(defaultWriter, "\t----------------")

	writer := tabwriter.NewWriter(os.Stdout, 0, 2, 10, ' ', tabwriter.Debug)
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
	return fmt.Sprintf("\t%s\t%s\t%s", s, v, stat)
}

func formatWorkerRow[T any](worker interfaces.Comparator[T]) string {
	s := strings.TrimSpace(worker.GetID().(string))
	workeri := worker.(*shared.ComparatorModule[T])
	s = colourize.Colourize(s, colourize.White)
	return fmt.Sprintf("%s\t%v\t\t%s\t", s, worker.TaskCount(), workeri.GetStatus())
}

func formatWorkerTitle(title []string) string {
	return fmt.Sprintf("%s\t  \t%v\t \t%s", title[0], title[1], title[2])
}

func printProgressBar(current, total int) {
	width := 50
	progress := float64(current) / float64(total)
	barWidth := int(progress * float64(width))

	bar := colourize.Colourize("["+string(repeat('#', barWidth))+string(repeat(' ', width-barWidth))+"]", colourize.Blue)
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
	msg := colourize.Colourize(fmt.Sprintf(" Comparator %s submitted a task.\n F: %s\n S: %s", p.AssignieeId, p.F, p.S), colourize.Green, colourize.Bold)
	fmt.Println(msg)
	newLine(1)
}

func printWorkerStatusTable[T any](workers []*(interfaces.Comparator[T])) {
	titleWriter := tabwriter.NewWriter(os.Stdout, 0, 5, 10, ' ', tabwriter.AlignRight)
	writer := tabwriter.NewWriter(os.Stdout, 10, 2, 10, ' ', tabwriter.AlignRight)

	// TODO: make this dynamic
	fmt.Fprintln(titleWriter, formatWorkerTitle([]string{"Comparator", "Task Count", "Status"}))
	fmt.Fprintln(titleWriter, formatWorkerTitle([]string{"---------", "----------", "------"}))

	for _, row := range workers {
		fmt.Fprintln(writer, formatWorkerRow(*row))
	}
	titleWriter.Flush()
	writer.Flush()
	newLine(2)

}
