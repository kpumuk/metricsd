package parser

import (
	"os"
	"testing"
	"metricsd/types"
)

type eventTest struct {
	buf     string
	results []testEntry
}

type testEntry struct {
	event *types.Event
	err   os.Error
}

var parseTests = []eventTest{
	// Valid events with single metric
	{"metric:10", []testEntry{
		{types.NewEvent("", "metric", 10), nil},
	}},
	{"metric:-1", []testEntry{
		{types.NewEvent("", "metric", -1), nil},
	}},
	{"group$metric:-1", []testEntry{
		{types.NewEvent("", "group$metric", -1), nil},
	}},
	{"group.metric:2", []testEntry{
		{types.NewEvent("", "group.metric", 2), nil},
	}},
	{"app01@metric:10", []testEntry{
		{types.NewEvent("app01", "metric", 10), nil},
	}},

	// Invalid events with single metric
	{":10", []testEntry{
		{nil, os.NewError("Metric name is empty (event=\":10\")")},
	}},
	{"metric!:10", []testEntry{
		{nil, os.NewError("Metric name is invalid: \"metric!\" (event=\"metric!:10\")")},
	}},
	{"src!@metric:10", []testEntry{
		{nil, os.NewError("Source is invalid: \"src!\" (event=\"src!@metric:10\")")},
	}},
	{"app01@metric", []testEntry{
		{nil, os.NewError("Event format is invalid (event=\"app01@metric\")")},
	}},
	{"app01@metric:hello", []testEntry{
		{nil, os.NewError("Metric value \"hello\" is invalid (event=\"app01@metric:hello\")")},
	}},

	// Valid events with multiple metrics
	{"metric1:10;metric2:20", []testEntry{
		{types.NewEvent("", "metric1", 10), nil},
		{types.NewEvent("", "metric2", 20), nil},
	}},
	{"app01@metric1:10;metric2:20", []testEntry{
		{types.NewEvent("app01", "metric1", 10), nil},
		{types.NewEvent("", "metric2", 20), nil},
	}},
	{"app01@metric1:10;app02@metric2:20", []testEntry{
		{types.NewEvent("app01", "metric1", 10), nil},
		{types.NewEvent("app02", "metric2", 20), nil},
	}},

	// Semi-valid events (multiple metrics, some are invalid)
	{"metric1:10;metric2:", []testEntry{
		{types.NewEvent("", "metric1", 10), nil},
		{nil, os.NewError("Metric value \"\" is invalid (event=\"metric1:10;metric2:\")")},
	}},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		var idx = 0
		count := Parse(test.buf, func(event *types.Event, err os.Error) {
			if idx == len(test.results) {
				t.Errorf("Unexpected event #%d: event=%q, err=%q (buf=%q, idx=%d)", idx, event, err, test.buf, idx)
				return
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
				if event == nil {
					t.Errorf("Expected event %q, got nil (buf=%q, idx=%d)", expected.event, test.buf, idx)
				}
				if expected.event == nil {
					t.Errorf("Expected nil, got event %q (buf=%q, idx=%d)", event, test.buf, idx)
				}

				if event != nil && expected.event != nil {
					// Test Event fields
					if event.Source != expected.event.Source {
						t.Errorf("Expected event source %q, got %q (buf=%q, idx=%d)", expected.event.Source, event.Source, test.buf, idx)
					}
					if event.Name != expected.event.Name {
						t.Errorf("Expected event name %q, got %q (buf=%q, idx=%d)", expected.event.Name, event.Name, test.buf, idx)
					}
					if event.Value != expected.event.Value {
						t.Errorf("Expected event value %q, got %q (buf=%q, idx=%d)", expected.event.Value, event.Name, test.buf, idx)
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

func BenchmarkParse(b *testing.B) {
	b.StopTimer()
	buf := "app01@group.metric:10;app02@group.metric:2;group.metric:2"
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Parse(buf, func(event *types.Event, err os.Error) {
			if err != nil {
				panic("Error occurred: " + err.String())
			}
		})
		b.SetBytes(int64(len(buf)))
	}

	b.StopTimer()
}
