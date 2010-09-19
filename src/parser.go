// The parser package implements gorrdpd protocol messages parsing.
//
// Basicly, message format is:
//     [source@]metric:value[;message]
// where source is the message source, metric and value - metric's name and value,
// and message is another message in the same format (you can send several metrics
// updates in the same package).
package parser

import (
    "fmt"
    "os"
    "strconv"
    "strings"
    "./types"
)

// Parse parses source buffer and invokes the given function, passing either parsed
// message or an error (when failed to parse) for each message in the source buffer
// (if there are several messages in the a bundle). Returns number of successfully
// processed messages.
//
// For example:
//     parser.Parse("app01@user_login:1;response_time:154;hello", func(msg *Message, err os.Error) {
//         fmt.Printf("Message=%v, Error=%v", msg, err)
//     })
// will invoke the given callback three times:
//     msg = &Message { Source: "app01", Name: "user_login",    Value: 1 },  err = nil
//     msg = &Message { Source: "",      Name: "response_time", Value: 154}, err = nil
//     msg = nil, err = os.Error (err.ToString() == "Message format is not valid: hello")
//
// Return value for this example will be 2.
func Parse(buf string, f func(message *types.Message, err os.Error)) int {
    // Number of successfully processed messages
    var count int
    // Process multiple metrics in a single message
    for _, msg := range strings.Split(buf, ";", -1) {
        var source, name, svalue string

        // Check if the message contains a source name
        if idx := strings.Index(msg, "@"); idx >= 0 {
            source = msg[:idx]
            msg    = msg[idx+1:]

            if !validateMetric(source) {
                f(nil, os.NewError(fmt.Sprintf("Source is invalid: %q (message=%q)", source, buf)))
                continue
            }
        }

        // Retrieve the metric name
        if idx := strings.Index(msg, ":"); idx >= 0 {
            name   = msg[:idx]
            svalue = msg[idx+1:]

            if !validateMetric(name) {
                f(nil, os.NewError(fmt.Sprintf("Metric name is invalid: %q (message=%q)", name, buf)))
                continue
            }
            if len(name) == 0 {
                f(nil, os.NewError(fmt.Sprintf("Metric name is empty (message=%q)", buf)))
                continue
            }
        } else {
            f(nil, os.NewError(fmt.Sprintf("Message format is invalid (message=%q)", buf)))
            continue
        }

        // Parse the value
        if value, error := strconv.Atoi(svalue); error != nil {
            f(nil, os.NewError(fmt.Sprintf("Metric value %q is invalid (message=%q)", svalue, buf)))
            continue
        } else {
            f(types.NewMessage(source, name, value), nil)
            count += 1
        }
    }
    return count
}

/***** Helper functions *******************************************************/

func validateMetric(name string) bool {
    for _, rune := range name {
        if rune < 0x80 {
            // Digits
            if '0' <= rune && rune <= '9' { continue }
            // Letters
            if 'a' <= rune && rune <= 'z' { continue }
            switch rune {
            // Special characters
            case '_', '-', '$', '.':
                    continue
            }
        }
        return false
    }
    return true
}
