package instrument

import (
	"github.com/IljaN/noi7gen/generator"
	"github.com/go-audio/audio"
)

type Instrument struct {
	SamplingRate int
	BitRate      int
	Duration     int
}

// Out is the fundamental building block to build signal flows. It takes an audio-buffer and modifies it.
type Out func(buf *audio.FloatBuffer)

func New(samplingRate int, bitRate int, duration int) *Instrument {
	p := &Instrument{
		SamplingRate: samplingRate,
		BitRate:      bitRate,
		Duration:     duration,
	}
	return p
}

// Master is the master output of the Instrument. Use this function to
// get the final sound
func (in *Instrument) Master(o Out) *audio.FloatBuffer {
	buf := in.EmptyBuf()
	o(buf)
	return buf
}

// EmptyBuf generates a suitable empty sample-buffer based on instruments config
func (in *Instrument) EmptyBuf() *audio.FloatBuffer {
	return &audio.FloatBuffer{
		Data:   make([]float64, in.Duration*2),
		Format: audio.FormatStereo48000,
	}
}

// NewOsc is a helper-function which creates an oscillator based on the instruments config
func (in *Instrument) NewOsc(shape generator.WaveFunc, hz float64) *generator.Osc {
	osc := generator.NewOsc(shape, hz, in.SamplingRate)
	osc.Amplitude = float64(audio.IntMaxSignedValue(in.BitRate))

	return osc
}

// OscOut returns a function which fills the audio-buffer with an oscillator/generator signal
func OscOut(osc *generator.Osc) func(buf *audio.FloatBuffer) {
	o := osc
	return func(buf *audio.FloatBuffer) {
		o.Fill(buf)
	}
}

// Chain creates a function which passes the audio-buffer trough multiple stages
// for transformation.
func Chain(stages ...Out) Out {
	return func(buf *audio.FloatBuffer) {
		for i := range stages {
			stages[i](buf)
		}
	}
}

// Mix symmetrically mixes multiple inputs
func (in *Instrument) Mix(inputs ...Out) Out {
	oscCount := float64(len(inputs))
	oscOuts := make([]*audio.FloatBuffer, 0, int(oscCount))

	// Pre allocate out buffers for each oscillator
	for i := 0; i != int(oscCount); i++ {
		oscOuts = append(oscOuts, in.EmptyBuf())
	}

	// Fill each oscOut with data from corresponding output
	return func(buf *audio.FloatBuffer) {
		for i := range inputs {
			inputs[i](oscOuts[i])
			for s := range buf.Data {
				buf.Data[s] += oscOuts[i].Data[s] / oscCount
			}
		}
	}
}

// Div returns a function which divides every sample of a given buffer by a static value
func Div(val float64) Out {
	return func(buf *audio.FloatBuffer) {
		for i := range buf.Data {
			if buf.Data[i] != 0.0 {
				buf.Data[i] /= val
			}
		}
	}
}
