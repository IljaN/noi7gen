package main

import (
	"fmt"
	"github.com/IljaN/noi7gen/encoding"
	"github.com/IljaN/noi7gen/filter"
	"github.com/IljaN/noi7gen/generator"
	inst "github.com/IljaN/noi7gen/instrument"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < 1; i++ {

		s := inst.New(48000, 16, 500000*2)

		var genRandOSCs = func() []inst.Out {
			oscNum := randInt(1, 3)
			oscs := make([]inst.Out, oscNum)
			for oN := range oscs {
				o := s.NewOsc(randomWf(), float64(randInt(60, 200)))
				o.SetAttackInMs(randInt(5, 3000))
				oscs[oN] = inst.OscOut(o)

			}

			return oscs

		}

		oscillators := genRandOSCs()

		out := s.Master(inst.Chain(s.Mix(oscillators...),
			filter.NewFlangerFilter(&filter.FlangerFilter{
				Time:    randFloats(0.1, 0.9999),
				Factor:  randFloats(0.0, 0.9999),
				LFORate: randFloats(0.003, 0.1),
			}),
			filter.NewLPF(&filter.LPF{
				Cutoff: randFloats(0.1, 0.999),
			}),
			filter.NewBitCrusher(&filter.BitCrusher{
				Factor: randFloats(0.1, 1.0),
			}),
			filter.NewDelayFilter(&filter.DelayFilter{
				LeftTime:      randFloats(0.3, 0.9),
				LeftFactor:    randFloats(0.1222, 0.7),
				LeftFeedback:  randFloats(0.1, 0.9),
				RightTime:     randFloats(0.1, 0.8),
				RightFactor:   randFloats(0.3, 0.8),
				RightFeedback: randFloats(0.2, 0.55),
			}),
		),
		)

		fileName := fmt.Sprintf("gen_%d.wav", i)
		if err := encoding.WriteWAV(fileName, out, 16); err != nil {
			log.Fatal(err)
		}

		ffplay("gen_0.wav", 48000)
	}

}

// plays and visualizes the generated sound with ffplay
func ffplay(fn string, sampleRate int) {
	ffplayExecutable, _ := exec.LookPath("ffplay")
	ffplayCmd := &exec.Cmd{
		Path:   ffplayExecutable,
		Args:   []string{ffplayExecutable, "-showmode", "1", "-ar", strconv.Itoa(sampleRate), fn},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := ffplayCmd.Start(); err != nil {
		panic(err)
	}

	ffplayCmd.Wait()
}

var wf = []generator.WaveFunc{
	generator.Sawtooth,
	generator.Sine,
	generator.Triangle,
	generator.Square,
	generator.WhiteNoise,
}

func randomWf() generator.WaveFunc {
	idx := randInt(0, len(wf))
	return wf[idx]
}

func randInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func randFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)

}
