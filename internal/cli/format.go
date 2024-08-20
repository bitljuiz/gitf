package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gitf/internal/statistics"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

func FormatData(stats statistics.SortedStatistics, mode string) error {
	switch mode {
	case "tabular":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintln(w, "Name\tLines\tCommits\tFiles\r")
		for _, stat := range stats {
			_, _ = fmt.Fprintf(w,
				strings.Join([]string{
					stat.Name,
					strconv.Itoa(stat.Lines),
					strconv.Itoa(stat.Commits),
					strconv.Itoa(stat.Files),
				}, "\t")+"%s", "\r\n")
		}
		if err := w.Flush(); err != nil {
			return err
		}
	case "csv":
		records := [][]string{
			{"Name", "Lines", "Commits", "Files"},
		}

		for _, stat := range stats {
			records = append(records,
				[]string{
					stat.Name,
					strconv.Itoa(stat.Lines),
					strconv.Itoa(stat.Commits),
					strconv.Itoa(stat.Files),
				},
			)
		}

		w := csv.NewWriter(os.Stdout)
		w.UseCRLF = true

		if err := w.WriteAll(records); err != nil {
			return err
		}

		w.Flush()

		if err := w.Error(); err != nil {
			return err
		}
	case "json":
		jsonStats, err := json.Marshal(&stats)

		if err != nil {
			return err
		}

		_, err = os.Stdout.Write(jsonStats)
		if err != nil {
			return err
		}
	case "json-lines":
		for _, stat := range stats {
			jsonStat, err := json.Marshal(&stat)

			if err != nil {
				return err
			}

			_, err = os.Stdout.Write(jsonStat)

			fmt.Printf("%s", "\r\n")
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown format %s.\nUse gitf --help for more information", mode)
	}
	return nil
}
