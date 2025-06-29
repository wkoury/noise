package brownnoise

import (
	"testing"
)

func TestNewBrownNoiseStreamer(t *testing.T) {
	tests := []struct {
		name     string
		damping  float64
		gain     float64
		stepSize float64
		wantErr  bool
	}{
		{
			name:     "valid parameters",
			damping:  0.9,
			gain:     0.5,
			stepSize: 0.02,
			wantErr:  false,
		},
		{
			name:     "invalid damping too low",
			damping:  0.0,
			gain:     0.5,
			stepSize: 0.02,
			wantErr:  true,
		},
		{
			name:     "invalid damping too high",
			damping:  1.1,
			gain:     0.5,
			stepSize: 0.02,
			wantErr:  true,
		},
		{
			name:     "invalid gain negative",
			damping:  0.9,
			gain:     -0.1,
			stepSize: 0.02,
			wantErr:  true,
		},
		{
			name:     "invalid step size negative",
			damping:  0.9,
			gain:     0.5,
			stepSize: -0.01,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			streamer, err := NewBrownNoiseStreamer(tt.damping, tt.gain, tt.stepSize)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewBrownNoiseStreamer() expected error but got none")
				}
				if streamer != nil {
					t.Errorf("NewBrownNoiseStreamer() expected nil streamer on error")
				}
			} else {
				if err != nil {
					t.Errorf("NewBrownNoiseStreamer() unexpected error: %v", err)
				}
				if streamer == nil {
					t.Errorf("NewBrownNoiseStreamer() expected valid streamer")
				}
			}
		})
	}
}

func TestBrownNoiseStreamer_Stream(t *testing.T) {
	streamer := &BrownNoiseStreamer{
		Accumulator: 0.0,
		Damping:     0.9,
		Gain:        0.5,
		StepSize:    0.02,
	}

	samples := make([][2]float64, 1024)
	n, ok := streamer.Stream(samples)

	if !ok {
		t.Errorf("Stream() returned ok=false, expected true")
	}
	if n != len(samples) {
		t.Errorf("Stream() returned n=%d, expected %d", n, len(samples))
	}

	// Check that samples are generated (not all zero)
	nonZero := false
	for _, sample := range samples {
		if sample[0] != 0.0 || sample[1] != 0.0 {
			nonZero = true
			break
		}
	}
	if !nonZero {
		t.Errorf("Stream() generated all zero samples")
	}
}

func TestBrownNoiseStreamer_Err(t *testing.T) {
	streamer := &BrownNoiseStreamer{}
	if err := streamer.Err(); err != nil {
		t.Errorf("Err() returned non-nil error: %v", err)
	}
}
