package parser

import (
    "os"
    "testing"
    "./types"
)

type messageTest struct {
    buf     string
    results []testEntry
}

type testEntry struct {
    message *types.Message
    err     os.Error
}

var parseTests = []messageTest{
    // Valid messages with single metric
    {"metric:10", []testEntry{
        {types.NewMessage("", "metric", 10), nil},
    }},
    {"metric:-1", []testEntry{
        {types.NewMessage("", "metric", -1), nil},
    }},
    {"group$metric:-1", []testEntry{
        {types.NewMessage("", "group$metric", -1), nil},
    }},
    {"app01@metric:10", []testEntry{
        {types.NewMessage("app01", "metric", 10), nil},
    }},

    // Invalid messages with single metric
    {":10", []testEntry{
        {nil, os.NewError("Metric name is empty (message=\":10\")")},
    }},
    {"metric!:10", []testEntry{
        {nil, os.NewError("Metric name is invalid: \"metric!\" (message=\"metric!:10\")")},
    }},
    {"src!@metric:10", []testEntry{
        {nil, os.NewError("Source is invalid: \"src!\" (message=\"src!@metric:10\")")},
    }},
    {"app01@metric", []testEntry{
        {nil, os.NewError("Message format is invalid (message=\"app01@metric\")")},
    }},
    {"app01@metric:hello", []testEntry{
        {nil, os.NewError("Metric value \"hello\" is invalid (message=\"app01@metric:hello\")")},
    }},

    // Valid messages with multiple metrics
    {"metric1:10;metric2:20", []testEntry{
        {types.NewMessage("", "metric1", 10), nil},
        {types.NewMessage("", "metric2", 20), nil},
    }},
    {"app01@metric1:10;metric2:20", []testEntry{
        {types.NewMessage("app01", "metric1", 10), nil},
        {types.NewMessage("", "metric2", 20), nil},
    }},
    {"app01@metric1:10;app02@metric2:20", []testEntry{
        {types.NewMessage("app01", "metric1", 10), nil},
        {types.NewMessage("app02", "metric2", 20), nil},
    }},

    // Semi-valid messages (multiple metrics, some are invalid)
    {"metric1:10;metric2:", []testEntry{
        {types.NewMessage("", "metric1", 10), nil},
        {nil, os.NewError("Metric value \"\" is invalid (message=\"metric1:10;metric2:\")")},
    }},
}

func TestParse(t *testing.T) {
    for _, test := range parseTests {
        var idx = 0
        count := Parse(test.buf, func(message *types.Message, err os.Error) {
            if idx == len(test.results) {
                t.Errorf("Unexpected message #%d: message=%q, err=%q (buf=%q, idx=%d)", idx, message, err, test.buf, idx)
            }

            expected := test.results[idx]
            // Test errors
            if err == nil && expected.err != nil {
                t.Errorf("Expected error %q, got no error (buf=%q, idx=%d)", expected.err, test.buf, idx)
            }
            if err != nil && expected.err == nil {
                t.Errorf("Expected no error, got error %q (buf=%q, idx=%d)", err, test.buf, idx)
            }
            if err != expected.err {
                t.Errorf("Expected error %q, got error %q (buf=%q, idx=%d)", expected.err, err, test.buf, idx)
            }
            if err == nil {
                if message == nil {
                    t.Errorf("Expected message %q, got nil (buf=%q, idx=%d)", expected.message, test.buf, idx)
                }
                if expected.message == nil {
                    t.Errorf("Expected nil, got message %q (buf=%q, idx=%d)", message, test.buf, idx)
                }

                if message != nil && expected.message != nil {
                    // Test Message fields
                    if message.Source != expected.message.Source {
                        t.Errorf("Expected message source %q, got %q (buf=%q, idx=%d)", expected.message.Source, message.Source, test.buf, idx)
                    }
                    if message.Name != expected.message.Name {
                        t.Errorf("Expected message name %q, got %q (buf=%q, idx=%d)", expected.message.Name, message.Name, test.buf, idx)
                    }
                    if message.Value != expected.message.Value {
                        t.Errorf("Expected message value %q, got %q (buf=%q, idx=%d)", expected.message.Value, message.Name, test.buf, idx)
                    }
                }
            }
            idx++
        })

        expectedCount := 0
        for _, result := range test.results {
            if result.err == nil {
                expectedCount++
            }
        }
        if count != expectedCount {
            t.Errorf("Expected to return %q, got %q (buf=%q)", len(test.results), count, test.buf)
        }
    }
}
