//go:build js && wasm
// +build js,wasm

package main

import (
	"noise/internal/brownnoise"
	"syscall/js"
)

var streamer *brownnoise.BrownNoiseStreamer

// initStreamer(damping, gain, stepSize)
func initStreamer(this js.Value, args []js.Value) any {
	streamer = &brownnoise.BrownNoiseStreamer{
		Accumulator: 0,
		Damping:     args[0].Float(),
		Gain:        args[1].Float(),
		StepSize:    args[2].Float(),
	}
	return nil
}

// nextSamples(frameCount) -> Float32Array of length frameCount*2
func nextSamples(this js.Value, args []js.Value) any {
	count := args[0].Int()
	buf := make([][2]float64, count)
	streamer.Stream(buf)
	// flatten to float32
	flat := make([]any, count*2)
	for i, pair := range buf {
		flat[i*2] = float32(pair[0])
		flat[i*2+1] = float32(pair[1])
	}
	return js.ValueOf(flat)
}

func registerCallbacks() {
	js.Global().Set("initStreamer", js.FuncOf(initStreamer))
	js.Global().Set("nextSamples", js.FuncOf(nextSamples))
}

func main() {
	registerCallbacks()
	// block forever
	select {}
}
