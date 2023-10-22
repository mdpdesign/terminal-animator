# Terminal Animator

A fun project to start learning Go.. Make some animations in terminal..

## Demo

Terminal Animator running in [Cool Retro Term](https://github.com/Swordfish90/cool-retro-term):

https://github.com/mdpdesign/terminal-animator/assets/5243838/04d0c2d3-7626-48da-a94e-182bd9794adb

## Docs

Usage:

```shell
terminal-animator -h

Usage of terminal-animator:
  -config string
        Configuration file path for animation
```

Example:

```shell
terminal-animator -config demo-animation-108.yaml
```

### Configuration documentation

Terminal Animator is configured by specifying a list of animation files:

```yaml
animationFiles:
  - demo-animation-width-108/001-boot.yaml
  - demo-animation-width-108/002-boot.yaml
  - demo-animation-width-108/003-load.yaml
  - demo-animation-width-108/004-welcome.yaml
  - (...)
```

Each file describes the animation type and additional settings for that animation

#### Animation files - supported fields:

- directives (`map[string]any`): Defines type and animation settings, see below
- frames (`[]string`): List of text frames of the animation

##### Directives:

- type (`string`): Animation type, currently supported animations: `clear-line`, `clear-screen`, `loop`, `printer`, `typewriter`
  - `clear-line`: Animation that renders frames a on single line, clearing any previous frame
  - `clear-screen`: Animation that redraws whole screen, clearing whole screen before each frame
  - `loop`: Animation that loops specified amount of times after rendering all frames
  - `printer`: Animation that renders each frame line by line, split by newline - `\n`
  - `typewriter`: Typewriter animation effect, either by character or by word
- maxDelay (`float32`): Random delay duration between rendering frames. Default `0.5`
- evenDelay (`bool`): Whether delay between rendering frames should be even
- endNewLines (`int`): Number of newlines to add after that particular animation
- loop (`int`): Number of how many times the animation will be repeated. Only applicable to `loop` animation
- word (`bool`): For `typewriter` animation, whether to split animation frame by word or by character (**default**)

Examples:

```yaml
directives:
  type: clear-line
  maxDelay: 0.5
  endNewLines: 2
frames:
  - 00010 KB OK
  - 00025 KB OK
  - 00048 KB OK
  - 00128 KB OK
  - 00168 KB OK
  - 00256 KB OK
  - 00640 KB OK
```

```yaml
directives:
  type: typewriter
  word: false # default, can be optional
  maxDelay: 0.2
  endNewLines: 0
frames:
  - '# A: load-welcome-message'
```

For more examples, check the demo animation

<p align=center>Coded with some ❤️ on my really old 💻 with Fedora</p>
