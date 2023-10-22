package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type IAnimation interface {
	Render()
}

type ClearLineAnimation struct {
	maxDelay    float32
	evenDelay   bool
	endNewLines int
	frames      []string
}

type ClearScreenAnimation struct {
	maxDelay    float32
	evenDelay   bool
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
	evenDelay   bool
	endNewLines int
	frames      []string
}

type TypewriterAnimation struct {
	word        bool
	maxDelay    float32
	evenDelay   bool
	endNewLines int
	frames      []string
}

func (a ClearLineAnimation) Render() {
	for _, frame := range a.frames {
		fmt.Printf("%s\r", clearText(frame))
		animationDelay(a.maxDelay, a.evenDelay)
	}
	endAnimation(a.endNewLines)
}

func (a ClearScreenAnimation) Render() {
	for _, frame := range a.frames {
		clearTerminal()
		fmt.Printf("%s", clearText(frame))
		animationDelay(a.maxDelay, a.evenDelay)
	}
	endAnimation(a.endNewLines)
}

func (a LoopAnimation) Render() {
	frames := a.frames

	for i := 0; i < a.loop; i++ {
		fmt.Printf("%s\r", frames[i%len(frames)])
		animationDelay(a.maxDelay, true)
	}

	endAnimation(a.endNewLines)
}

func (a PrinterAnimation) Render() {
	for _, frame := range a.frames {
		frameParts := strings.Split(frame, "\n")

		for _, line := range frameParts {
			line = clearText(line)
			fmt.Printf("%s\n", line)
			animationDelay(a.maxDelay, a.evenDelay)
		}
	}
	endAnimation(a.endNewLines)
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
			animationDelay(a.maxDelay, a.evenDelay)
		}

		fmt.Println("")
	}
	endAnimation(a.endNewLines)
}

func NewAnimation(config *FrameConfig) (IAnimation, error) {
	var a IAnimation

	switch config.Directives.Type {
	case "typewriter":
		a = &TypewriterAnimation{
			config.Directives.Word,
			config.Directives.MaxDelay,
			config.Directives.EvenDelay,
			config.Directives.EndNewLines,
			config.Frames,
		}
	case "printer":
		a = &PrinterAnimation{
			config.Directives.MaxDelay,
			config.Directives.EvenDelay,
			config.Directives.EndNewLines,
			config.Frames,
		}
	case "loop":
		a = &LoopAnimation{
			config.Directives.Loop,
			config.Directives.MaxDelay,
			config.Directives.EndNewLines,
			config.Frames,
		}
	case "clear-line":
		a = &ClearLineAnimation{
			config.Directives.MaxDelay,
			config.Directives.EvenDelay,
			config.Directives.EndNewLines,
			config.Frames,
		}
	case "clear-screen":
		a = &ClearScreenAnimation{
			config.Directives.MaxDelay,
			config.Directives.EvenDelay,
			config.Directives.EndNewLines,
			config.Frames,
		}
	default:
		return nil, fmt.Errorf(
			"animation type '%s' not supported, allowed types: %s",
			config.Directives.Type,
			getAnimationtypes(),
		)
	}

	return a, nil
}

func getAnimationtypes() []string {
	return []string{"typewriter", "printer", "loop", "clear-line", "clear-screen"}
}

func animationDelay(maxDelay float32, even bool) {
	if even {
		time.Sleep(time.Duration(maxDelay * float32(time.Second)))
	} else {
		r := 0 + rand.Float32()*maxDelay
		time.Sleep(time.Duration(r * float32(time.Second)))
	}
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
