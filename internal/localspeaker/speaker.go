package localspeaker

import (
	"noise/internal/brownnoise"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

func Play() {
	// Sample rate for “CD-quality” audio
	const sampleRate = 44100

	// Initialize the speaker
	speaker.Init(beep.SampleRate(sampleRate), sampleRate/10)

	// Create our infinite Brown noise streamer with a damping factor
	brown := &brownnoise.BrownNoiseStreamer{
		Accumulator: 0.0,
		Damping:     0.90,
		Gain:        0.5,  // reduce volume to soften
		StepSize:    0.02, // reduce step size to further smooth the noise
	}

	// Play it indefinitely
	speaker.Play(brown)

	// Block forever (or until you kill the program)
	select {}
}
