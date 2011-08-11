// The parser package implements MetricsD protocol events parsing.
//
// Basicly, event format is:
//     [source@]metric:value[;event]
// where source is the event source, metric and value - metric's name and value,
// and event is another event in the same format (you can send several metrics
// updates in the same package).
package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"metricsd/types"
)

// Parse parses source buffer and invokes the given function, passing either parsed
// event or an error (when failed to parse) for each event in the source buffer
// (if there are several events in the a bundle). Returns number of successfully
// processed events.
//
// For example:
//     parser.Parse("app01@user_login:1;response_time:154;hello", func(msg *event, err os.Error) {
//         fmt.Printf("event=%v, Error=%v", msg, err)
//     })
// will invoke the given callback three times:
//     msg = &Event { Source: "app01", Name: "user_login",    Value: 1 },  err = nil
//     msg = &Event { Source: "",      Name: "response_time", Value: 154}, err = nil
//     msg = nil, err = os.Error (err.ToString() == "Event format is invalid: hello")
//
// Return value for this example will be 2.
func Parse(buf string, f func(event *types.Event, err os.Error)) int {
	// Number of successfully processed events
	var count int
	// Process multiple metrics in a single event
	var msg string
	for str := buf; str != ""; {
		if idx := strings.Index(str, ";"); idx >= 0 {
			msg, str = str[:idx], str[idx+1:]
		} else {
			msg, str = str, str[:0]
		}

		var source, name, svalue string

		// Check if the event contains a source name
		if idx := strings.Index(msg, "@"); idx >= 0 {
			source, msg = msg[:idx], msg[idx+1:]

			if !validateMetric(source) {
				f(nil, os.NewError(fmt.Sprintf("Source is invalid: %q (event=%q)", source, buf)))
				continue
			}
		}

		// Retrieve the metric name
		if idx := strings.Index(msg, ":"); idx >= 0 {
			name, svalue = msg[:idx], msg[idx+1:]

			if !validateMetric(name) {
				f(nil, os.NewError(fmt.Sprintf("Metric name is invalid: %q (event=%q)", name, buf)))
				continue
			}
			if len(name) == 0 {
				f(nil, os.NewError(fmt.Sprintf("Metric name is empty (event=%q)", buf)))
				continue
			}
		} else {
			f(nil, os.NewError(fmt.Sprintf("Event format is invalid (event=%q)", buf)))
			continue
		}

		// Parse the value
		if value, error := strconv.Atoi(svalue); error != nil {
			f(nil, os.NewError(fmt.Sprintf("Metric value %q is invalid (event=%q)", svalue, buf)))
			continue
		} else {
			f(types.NewEvent(source, name, value), nil)
			count += 1
		}
	}
	return count
}

/***** Helper functions *******************************************************/

func validateMetric(name string) bool {
	for _, rune := range name {
		if rune > 0x7F {
			return false
		}

		// Digits and Letters
		if ('0' <= rune && rune <= '9') || ('a' <= rune && rune <= 'z') || ('A' <= rune && rune <= 'Z') {
			continue
		}
		// Special characters
		switch rune {
		case '_', '-', '$', '.':
			continue
		default:
			return false
		}
	}
	return true
}
