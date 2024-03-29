package main

import (
	"flag"
	"log"
	"os"

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
		EvenDelay   bool    `yaml:"evenDelay,omitempty"`
		EndNewLines int     `yaml:"endNewLines,omitempty"`
	}
	Frames []string `yaml:"frames"`
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

		if fCfg.Directives.MaxDelay == 0 {
			fCfg.Directives.MaxDelay = 0.5
		}

		a, err := NewAnimation(&fCfg)
		if err != nil {
			log.Fatalf("error creating animation for: %s, %s", file, err)
		}
		a.Render()
	}
}
