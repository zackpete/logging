package logging

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type Console struct{}

const timeFormat = "15:04:05.000"

func (this Console) Handle(event *Event) {
	var (
		w io.Writer
		l string
	)

	switch event.Level {
	case Error:
		w = os.Stderr
		l = "ERR"
	case Warn:
		w = os.Stderr
		l = "WRN"
	case Info:
		w = os.Stdout
		l = "INF"
	case Debug:
		w = os.Stdout
		l = "DBG"
	}

	_, _ = fmt.Fprintf(w, "%s [%s] %s%s\n",
		event.Time.Format(timeFormat), l, event.Message, formatData(event.Context))
}

func formatData(data Data) string {
	if len(data) == 0 {
		return ""
	}

	values := make([]string, 0, len(data))
	for k := range data {
		values = append(values, k)
	}

	sort.Strings(values)
	for i, k := range values {
		values[i] = fmt.Sprintf("%s = %v", k, data[k])
	}

	return " | " + strings.Join(values, ", ")
}
