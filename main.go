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
		Type        string  `yaml:"type"`
		Loop        int     `yaml:"loop,omitempty"`
		Word        bool    `yaml:"word,omitempty"`
		MaxDelay    float32 `yaml:"maxDelay,omitempty"`
		EndNewLines int     `yaml:"endNewLines,omitempty"`
	}
	Frames []string `yaml:"frames"`
}

func getAnimationtypes() []string {
	return []string{"clear-line", "clear-screen", "loop", "printer"}
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

		animationType := fCfg.Directives.Type
		loop := fCfg.Directives.Loop
		typeWriterWordSplit := fCfg.Directives.Word
		maxDelay := fCfg.Directives.MaxDelay
		endNewLines := fCfg.Directives.EndNewLines

		if maxDelay == 0 {
			maxDelay = 0.5
		}

		for _, frame := range fCfg.Frames {
			switch animationType {
			case "typewriter":

				var format string
				var frameParts []string

				if typeWriterWordSplit {
					frameParts = strings.Split(frame, " ")
					format = "%s "
				} else {
					frameParts = strings.Split(frame, "")
					format = "%s"
				}

				for _, el := range frameParts {
					fmt.Printf(format, el)
					randomDelay(maxDelay)
				}

				fmt.Println("")
			case "printer":
				frameParts := strings.Split(frame, "\n")

				for _, line := range frameParts {
					line = clearText(line)
					fmt.Printf("%s\n", line)
					randomDelay(maxDelay)
				}
			case "loop":
				frames := fCfg.Frames

				for i := 0; i < loop; i++ {
					fmt.Printf("%s\r", frames[i%len(frames)])
					time.Sleep(time.Duration(maxDelay * float32(time.Second)))
				}
			case "clear-line":
				fmt.Printf("%s\r", clearText(frame))
				randomDelay(maxDelay)
			case "clear-screen":
				clearTerminal()

				fmt.Printf("%s", clearText(frame))
				randomDelay(maxDelay)
			default:
				clearTerminal()

				log.Fatalf(
					"error animation type '%s' not supported, allowed types: %s",
					fCfg.Directives.Type,
					getAnimationtypes(),
				)
			}
		}

		if endNewLines > 0 {
			newLinesToPrint := strings.Repeat("\n", endNewLines)
			fmt.Print(newLinesToPrint)
		}
	}
}
