package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type YamlConfig struct {
	AnimationFiles []string `yaml:"animationFiles"`
}

type FrameConfig struct {
	Directives struct {
		Clear       bool    `yaml:"clear,omitempty"`
		FullRedraw  bool    `yaml:"fullRedraw,omitempty"`
		SplitByLine bool    `yaml:"splitByLine,omitempty"`
		Loop        int     `yaml:"loop,omitempty"`
		MaxDelay    float32 `yaml:"maxDelay,omitempty"`
		EndNewLines int     `yaml:"endNewLines,omitempty"`
	}
	Frames []string `yaml:"frames"`
}

func clearText(s string) string {
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "'", "")

	return s
}

func randomDelay(maxDelay float32) {
	r := 0 + rand.Float32()*maxDelay
	time.Sleep(time.Duration(r * float32(time.Second)))
}

func clearTerminal() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func main() {
	fileList := YamlConfig{}

	configFile := flag.String("config", "", "Configuration file path for animation")
	flag.Parse()

	if *configFile == "" {
		flag.CommandLine.Usage()
		os.Exit(2)
	}

	data, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("error occurred while reading file: %s, %v", *configFile, err)
	}

	err = yaml.Unmarshal([]byte(data), &fileList)
	if err != nil {
		log.Fatalf("error parsing Yaml file: %v", err)
	}

	clearTerminal()

	for _, file := range fileList.AnimationFiles {
		fCfg := FrameConfig{}

		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("error occurred while reading file: %s, %v", file, err)
		}

		err = yaml.Unmarshal([]byte(data), &fCfg)
		if err != nil {
			log.Fatalf("error parsing Yaml file: %v", err)
		}

		clear := fCfg.Directives.Clear
		fullRedraw := fCfg.Directives.FullRedraw
		loop := fCfg.Directives.Loop
		splitByLine := fCfg.Directives.SplitByLine
		maxDelay := fCfg.Directives.MaxDelay
		endNewLines := fCfg.Directives.EndNewLines

		end := "\n"
		if clear {
			end = "\r"
		}

		if maxDelay == 0 {
			maxDelay = 0.5
		}

		for _, frame := range fCfg.Frames {
			if loop > 0 {
				frames := fCfg.Frames

				for i := 0; i < loop; i++ {
					fmt.Printf("%s\r", frames[i%len(frames)])
					time.Sleep(time.Duration(maxDelay * float32(time.Second)))
				}
			} else if splitByLine {
				frameParts := strings.Split(frame, "\n")

				for _, line := range frameParts {
					line = clearText(line)
					fmt.Printf("%s%s", line, end)
					randomDelay(maxDelay)
				}
			} else {
				if fullRedraw {
					clearTerminal()
				}

				fmt.Printf("%s%s", clearText(frame), end)
				randomDelay(maxDelay)
			}
		}

		if endNewLines > 0 {
			newLinesToPrint := strings.Repeat("\n", endNewLines)
			fmt.Print(newLinesToPrint)
		}
	}
}
