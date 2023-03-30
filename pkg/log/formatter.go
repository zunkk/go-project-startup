package log

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	DisableCallerField                = "__format_disable_caller"
	DisableTrimMessageField           = "__format_disable_trim_message"
	OnlyWriteMsgWithoutFormatterField = "__format_only_write_msg_without_formatter"
)

var internalFields = map[string]struct{}{
	DisableCallerField:                {},
	DisableTrimMessageField:           {},
	OnlyWriteMsgWithoutFormatterField: {},
}

// Formatter - logrus formatter, implements logrus.Formatter
type Formatter struct {
	FirstFieldsOrder []string
	LastFieldsOrder  []string
	TimestampFormat  string
}

// Format an log entry
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	if _, ok := entry.Data[OnlyWriteMsgWithoutFormatterField]; ok {
		return []byte(entry.Message), nil
	}

	levelColor := getColorByLevel(entry.Level)

	disableCaller := false
	disableTrimMessage := false
	if _, ok := entry.Data[DisableCallerField]; ok {
		disableCaller = true
	}
	if _, ok := entry.Data[DisableTrimMessageField]; ok {
		disableTrimMessage = true
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}

	// output buffer
	b := &bytes.Buffer{}

	// write level
	level := strings.ToUpper(entry.Level.String())
	_, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m ", levelColor, level[:4])

	// write time
	b.WriteString(fmt.Sprintf("\x1b[%dm", colorGreen))
	b.WriteString("[")
	b.WriteString(entry.Time.Format(timestampFormat))
	b.WriteString("] ")
	b.WriteString("\x1b[0m")

	// write msg
	msg := entry.Message
	if !disableTrimMessage {
		msg = strings.TrimSpace(msg)
	}
	b.WriteString(msg)

	b.WriteString(fmt.Sprintf("\x1b[%dm", levelColor))
	// write fields
	if len(f.FirstFieldsOrder) == 0 && len(f.LastFieldsOrder) == 0 {
		f.writeFields(b, entry)
	} else {
		f.writeOrderedFields(b, entry)
	}
	b.WriteString("\x1b[0m")

	if !disableCaller {
		b.WriteString(fmt.Sprintf("\x1b[%dm", colorBlue))
		f.writeCaller(b, entry)
		b.WriteString("\x1b[0m")
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *Formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		f := ""
		ps := strings.Split(entry.Caller.File, "/")
		if len(ps) >= 2 {
			packName, _, _ := strings.Cut(ps[len(ps)-2], "@")
			f += packName + "/"
		}
		if len(ps) >= 1 {
			f += ps[len(ps)-1]
		}
		_, _ = fmt.Fprintf(b, " %s:%d", f, entry.Caller.Line)
	}
}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FirstFieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	lastFields := new(bytes.Buffer)
	for _, field := range f.LastFieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(lastFields, entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if !foundFieldsMap[field] {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}

	if lastFields.Len() != 0 {
		b.Write(lastFields.Bytes())
	}
}

func (f *Formatter) writeField(b io.Writer, entry *logrus.Entry, field string) {
	if _, ok := internalFields[field]; ok {
		return
	}
	_, _ = fmt.Fprintf(b, " [%s:%v]", field, entry.Data[field])
}

const (
	colorRed    = 31
	colorGreen  = 32
	colorYellow = 33
	colorBlue   = 34
	colorCyan   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorCyan
	}
}
