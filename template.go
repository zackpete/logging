package logging

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"time"
)

var (
	DefaultDateFormat = "2006-01-02T15:04:05.000"
)

var (
	r1       = regexp.MustCompile("{([^}]+)}")
	r2       = regexp.MustCompile("^[a-z][a-z0-9_]*(|.+)$")
	escaped0 = []byte{0, 0}
	escaped1 = []byte{1, 1}
	braceL   = []byte{'{'}
	braceR   = []byte{'}'}
)

func convert(template string, args []interface{}) (string, Data) {
	templateBytes := []byte(template)
	escapedL, escapedR := escape(templateBytes)

	templateMatches := r1.FindAllSubmatchIndex(templateBytes, -1)

	data := make(Data)
	formatBuffer := new(bytes.Buffer)
	messageArgs := make([]interface{}, 0, len(args))

	a := 0 // arg index
	n := 0 // template byte index
	for _, match := range templateMatches {
		formatBuffer.Write(templateBytes[n : match[2]-1])

		// for '{example|%s}' -> 'example|%s'
		propTemplate := templateBytes[match[2]:match[3]]
		propMatch := r2.FindSubmatchIndex(propTemplate)

		var (
			key string
			arg interface{}
		)

		if propMatch == nil { // {...} contains non-matching characters
			formatBuffer.Write(templateBytes[match[0]:match[1]])
			goto advance
		}

		key = string(propTemplate[:propMatch[2]])

		if a >= len(args) { // more properties than arguments
			_, _ = fmt.Fprintf(formatBuffer, "{%s:MISSING}", key)
			args = append(args, nil)
			goto advance
		}

		arg = resolveArg(args[a])

		if propMatch[2] == propMatch[3] { // default format
			formatArg(formatBuffer, arg)
		} else { // custom format, write '%s' for 'example|%s'
			formatBuffer.Write(propTemplate[propMatch[2]+1 : propMatch[3]])
			messageArgs = append(messageArgs, arg)
		}

		data[key] = arg
		a++

	advance:

		n = match[3] + 1
	}

	for ; a < len(args); a++ { // if extra args
		data[fmt.Sprintf("_%d", a)] = resolveArg(args[a])
	}

	formatBuffer.Write(templateBytes[n:])
	formatBytes := formatBuffer.Bytes()

	if escapedL {
		formatBytes = bytes.ReplaceAll(formatBytes, escaped0, braceL)
	}
	if escapedR {
		formatBytes = bytes.ReplaceAll(formatBytes, escaped1, braceR)
	}

	return fmt.Sprintf(string(formatBytes), messageArgs...), data
}

func resolveArg(arg interface{}) interface{} {
	if v := reflect.ValueOf(arg); v.Kind() == reflect.Func {
		if v.Type().NumIn() != 0 {
			return arg
		}

		result := v.Call(nil)

		if len(result) > 0 {
			return result[0].Interface()
		} else {
			return arg
		}
	} else {
		return arg
	}
}

func formatArg(w io.Writer, arg interface{}) {
	switch v := arg.(type) {
	case time.Time:
		if MessagesUTC {
			v = v.UTC()
		}
		_, _ = fmt.Fprint(w, v.Format(DefaultDateFormat))

	default:
		_, _ = fmt.Fprintf(w, "%v", arg)
	}
}

func escape(bytes []byte) (escapedL bool, escapedR bool) {
	var prev byte
	for i, b := range bytes {
		switch b {
		case '{':
			if prev == '{' {
				escapedL = true
				bytes[i-1] = 0
				bytes[i] = 0
				prev = 0
			} else {
				prev = '{'
			}
		case '}':
			if prev == '}' {
				escapedR = true
				bytes[i-1] = 1
				bytes[i] = 1
				prev = 0
			} else {
				prev = '}'
			}
		default:
			prev = 0
		}
	}

	return
}
