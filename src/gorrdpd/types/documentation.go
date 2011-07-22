// The types package implements various types needed by gorrdpd to store
// data.
//
// Incoming metrics are stored in instances of Message struct and hold
// information about source, metric's name, and value.
//
// The most interesting part is the structure holding metrics after they
// was parsed, but before aggregated and saved into the RRD. Basically,
// it could be represented in following way:
//
//     Slices -> []Slice -> []SampleSet
//
// Slices is a root structure, where slices being stored. Slice basically is
// an list of metrics with their values for a given period of time, which
// depends on SliceInterval config value. For example, if SliceInterval is 10,
// period beginnings will be the multiples of 10 (e.g., 1283473320-1283473329).
//
// SampleSet is a set of the specific metric values for a given source; the same
// metric for different sources will be stored in different sample sets. There is
// a special sample set, which source name is "all", where all values from all
// sources for a given metric are stored (useful to build summary stats for a metric).
//
// There are two primary tasks could be done using this package:
//
// 1. Add a message to Slices:
//     // Somewhere in the beginning, there usually only on Slices instance
//     slices := types.NewSlices(10)
//     // When you receive message
//     msg := NewMessage("app01", "user_login", 1)
//     slices.Add(msg)
// 2. Retrieving "closed" slices (or sample sets) to process them and store in some DB:
//     // retrieve closed slices
//     var closedSlices *vector.Vector = slices.ExtractClosedSlices(false)
//     // or retrieve all sample set for all closed slices in a single list
//     var closedSampleSets *vector.Vector = slices.ExtractClosedSampleSets(false)
//

package types