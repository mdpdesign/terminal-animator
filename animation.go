package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Animation interface {
	Render()
}

type ClearLineAnimation struct {
	maxDelay    float32
	endNewLines int
	frames      []string
}

type ClearScreenAnimation struct {
	maxDelay    float32
	endNewLines int
	frames      []string
}

type LoopAnimation struct {
	loop        int
	maxDelay    float32
	endNewLines int
	frames      []string
}

type PrinterAnimation struct {
	maxDelay    float32
	endNewLines int
	frames      []string
}

type TypewriterAnimation struct {
	word        bool
	maxDelay    float32
	endNewLines int
	frames      []string
}

func NewClearLineAnimation(maxDelay float32, endNewLines int, frames []string) *ClearLineAnimation {
	return &ClearLineAnimation{
		maxDelay:    maxDelay,
		endNewLines: endNewLines,
		frames:      frames,
	}
}
func (a ClearLineAnimation) Render() {
	for _, frame := range a.frames {
		fmt.Printf("%s\r", clearText(frame))
		randomDelay(a.maxDelay)
	}
	endAnimation(a.endNewLines)
}

func NewClearScreenAnimation(maxDelay float32, endNewLines int, frames []string) *ClearScreenAnimation {
	return &ClearScreenAnimation{
		maxDelay:    maxDelay,
		endNewLines: endNewLines,
		frames:      frames,
	}
}
func (a ClearScreenAnimation) Render() {
	for _, frame := range a.frames {
		clearTerminal()
		fmt.Printf("%s", clearText(frame))
		randomDelay(a.maxDelay)
	}
	endAnimation(a.endNewLines)
}

func NewLoopAnimation(loop int, endNewLines int, frames []string) *LoopAnimation {
	return &LoopAnimation{
		loop:        loop,
		maxDelay:    0.5,
		endNewLines: endNewLines,
		frames:      frames,
	}
}
func (a LoopAnimation) Render() {
	frames := a.frames

	for i := 0; i < a.loop; i++ {
		fmt.Printf("%s\r", frames[i%len(frames)])
		time.Sleep(time.Duration(a.maxDelay * float32(time.Second)))
	}

	endAnimation(a.endNewLines)
}

func NewPrinterAnimation(maxDelay float32, endNewLines int, frames []string) *PrinterAnimation {
	return &PrinterAnimation{
		maxDelay:    maxDelay,
		endNewLines: endNewLines,
		frames:      frames,
	}
}
func (a PrinterAnimation) Render() {
	for _, frame := range a.frames {
		frameParts := strings.Split(frame, "\n")

		for _, line := range frameParts {
			line = clearText(line)
			fmt.Printf("%s\n", line)
			randomDelay(a.maxDelay)
		}
	}
	endAnimation(a.endNewLines)
}

func NewTypewriterAnimation(word bool, maxDelay float32, endNewLines int, frames []string) *TypewriterAnimation {
	return &TypewriterAnimation{
		word:        word,
		maxDelay:    maxDelay,
		endNewLines: endNewLines,
		frames:      frames,
	}
}
func (a TypewriterAnimation) Render() {
	for _, frame := range a.frames {
		var format string
		var frameParts []string

		if a.word {
			frameParts = strings.Split(frame, " ")
			format = "%s "
		} else {
			frameParts = strings.Split(frame, "")
			format = "%s"
		}

		for _, el := range frameParts {
			fmt.Printf(format, el)
			randomDelay(a.maxDelay)
		}

		fmt.Println("")
	}
	endAnimation(a.endNewLines)
}

func MakeAnimation(config FrameConfig) Animation {
	var a Animation

	switch config.Directives.Type {
	case "typewriter":
		a = NewTypewriterAnimation(
			config.Directives.Word,
			config.Directives.MaxDelay,
			config.Directives.EndNewLines,
			config.Frames,
		)
	case "printer":
		a = NewPrinterAnimation(
			config.Directives.MaxDelay,
			config.Directives.EndNewLines,
			config.Frames,
		)
	case "loop":
		a = NewLoopAnimation(
			config.Directives.Loop,
			config.Directives.EndNewLines,
			config.Frames,
		)
	case "clear-line":
		a = NewClearLineAnimation(
			config.Directives.MaxDelay,
			config.Directives.EndNewLines,
			config.Frames,
		)
	case "clear-screen":
		a = NewClearScreenAnimation(
			config.Directives.MaxDelay,
			config.Directives.EndNewLines,
			config.Frames,
		)
	default:
		clearTerminal()

		log.Fatalf(
			"error animation type '%s' not supported, allowed types: %s",
			config.Directives.Type,
			getAnimationtypes(),
		)
	}

	return a
}

func getAnimationtypes() []string {
	return []string{"typewriter", "printer", "loop", "clear-line", "clear-screen"}
}

func randomDelay(maxDelay float32) {
	r := 0 + rand.Float32()*maxDelay
	time.Sleep(time.Duration(r * float32(time.Second)))
}

func clearText(s string) string {
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "'", "")

	return s
}

func clearTerminal() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func endAnimation(endNewLines int) {
	if endNewLines > 0 {
		newLinesToPrint := strings.Repeat("\n", endNewLines)
		fmt.Print(newLinesToPrint)
	}
}
