package main

import (
	"math/rand"
	"time"
)

type FrameConfig struct {
	Directives struct {
		Clear       bool    `yaml:"clear,omitempty"`
		FullRedraw  bool    `yaml:"fullRedraw,omitempty"`
		SplitByLine bool    `yaml:"splitByLine,omitempty"`
		Loop        int     `yaml:"loop,omitempty"`
		MaxDelay    float32 `yaml:"maxDelay,omitempty"`
		EndNewLines int     `yaml:"endNewLines,omitempty"`
	}
	Frames []string
}

func randomDelay(maxDelay float32) {
	r := 0 + rand.Float32()*maxDelay
	time.Sleep(time.Duration(r * float32(time.Second)))
}

func main() {

}
