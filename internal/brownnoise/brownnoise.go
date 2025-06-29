package brownnoise

import (
	"fmt"
	"math/rand"
)

// BrownNoiseStreamer implements beep.Streamer for brown (Brownian) noise.
type BrownNoiseStreamer struct {
	// Current "position" of the random walk
	Accumulator float64
	// Multiplicative factor to keep the walk from wandering too far
	Damping float64
	// Volume control to soften the noise
	Gain float64
	// Step size to control the smoothness of the noise
	StepSize float64
}

// NewBrownNoiseStreamer creates a new BrownNoiseStreamer with validation
func NewBrownNoiseStreamer(damping, gain, stepSize float64) (*BrownNoiseStreamer, error) {
	if damping <= 0 || damping > 1 {
		return nil, fmt.Errorf("damping must be between 0 and 1, got %f", damping)
	}
	if gain < 0 || gain > 1 {
		return nil, fmt.Errorf("gain must be between 0 and 1, got %f", gain)
	}
	if stepSize <= 0 {
		return nil, fmt.Errorf("stepSize must be positive, got %f", stepSize)
	}

	return &BrownNoiseStreamer{
		Accumulator: 0.0,
		Damping:     damping,
		Gain:        gain,
		StepSize:    stepSize,
	}, nil
}

// Stream fills the provided sample slice with brown noise samples.
func (b *BrownNoiseStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		// Generate uniform white noise in -1..+1, scaled by stepSize
		white := (rand.Float64()*2 - 1) * b.StepSize

		// Integrate with damping
		b.Accumulator = b.Accumulator*b.Damping + white

		// Optionally clamp to prevent large excursions
		if b.Accumulator > 1.0 {
			b.Accumulator = 1.0
		} else if b.Accumulator < -1.0 {
			b.Accumulator = -1.0
		}

		// Same signal on left and right channels with gain applied
		samples[i][0] = b.Accumulator * b.Gain
		samples[i][1] = b.Accumulator * b.Gain
	}
	return len(samples), true // return all samples and say “keep going”
}

// Err is required by beep.Streamer but we have no errors here.
func (b *BrownNoiseStreamer) Err() error {
	return nil
}
