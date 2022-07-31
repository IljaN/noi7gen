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
			oscNum := randInt(1, 7)
			oscs := make([]inst.Out, oscNum)
			atk := randInt(10, 500)
			for oN := range oscs {
				o := s.NewOsc(randomWf(), float64(randInt(60, 300)))
				o.SetAttackInMs(atk)
				oscs[oN] = inst.OscOut(o)

			}

			return oscs

		}

		oscillators := genRandOSCs()

		out := s.Master(inst.Chain(s.Mix(oscillators...),
			filter.NewFlangerFilter(&filter.FlangerFilter{
				Time:    randFloat(0.1, 0.9999),
				Factor:  randFloat(0.0, 0.9999),
				LFORate: randFloat(0.003, 0.1),
			}),
			filter.NewLPF(&filter.LPF{
				Cutoff: randFloat(0.1, 0.999),
			}),

			filter.NewBitCrusher(&filter.BitCrusher{
				Factor: randFloat(0.1, 0.999),
			}),

			filter.NewSimpleConvolutionFilter(&filter.SimpleConvolutionFilter{Coefficients: nRandFloats(200, 0.003, 0.005)}),
			filter.NewDelayFilter(&filter.DelayFilter{
				LeftTime:      randFloat(0.3, 0.9),
				LeftFactor:    randFloat(0.1222, 0.7),
				LeftFeedback:  randFloat(0.132, 0.9),
				RightTime:     randFloat(0.1, 0.8),
				RightFactor:   randFloat(0.3, 0.8),
				RightFeedback: randFloat(0.2, 0.55),
			}),
		),
		)

		fileName := fmt.Sprintf("gen_%d.wav", i)

		if err := encoding.WriteWAV(fileName, out, 16); err != nil {
			log.Fatal(err)
		}

		fmt.Println(fileName + " generated")
		ffplay(fileName, 48000)
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
	//generator.WhiteNoise,
}

func randomWf() generator.WaveFunc {
	idx := randInt(0, len(wf))
	return wf[idx]
}

func randInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func nRandFloats(n int, min, max float64) []float64 {
	floats := make([]float64, n)
	for i := range floats {
		floats[i] = randFloat(min, max)
	}

	return floats
}
