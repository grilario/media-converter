package command

import (
	"testing"
)

// func TestVideoOut_AddStream(t *testing.T) {
// 	type fields struct {
// 		Filename string
// 		Streams  []StreamOut
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			out := &Output{
// 				Filename: tt.fields.Filename,
// 				Streams:  tt.fields.Streams,
// 			}
// 			out.AddStream()
// 		})
// 	}
// }

func TestOutput_GenerateCommand(t *testing.T) {
	out := &Output{}

	streamIn := &Stream{
		Name:  "video",
		Codec: "h264",

		index:      1,
		inputIndex: 0,
	}
	out.Streams = append(out.Streams, StreamOut{ref: streamIn, Stream: Stream{Codec: "h265", index: 2}})

	got := out.GenerateCommand()

	t.Errorf("%v", got)

	// type fields struct {
	// 	Filename string
	// 	Streams  []StreamOut
	// }
	// tests := []struct {
	// 	name   string
	// 	fields fields
	// 	want   string
	// }{
	// 	// TODO: Add test cases.
	//    {
	//      "Simple test",
	//      fields{
	//
	//      }
	//    }
	// }
	//
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		out := &Output{
	// 			Filename: tt.fields.Filename,
	// 			Streams:  tt.fields.Streams,
	// 		}
	// 		if got := out.GenerateCommand(); got != tt.want {
	// 			t.Errorf("Output.GenerateCommand() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}
