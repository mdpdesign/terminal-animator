package main

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func Test_clearText(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Clear double quotes",
			args: args{
				s: "String with \"double\" quotes",
			},
			want: "String with double quotes",
		},
		{
			name: "Clear single quotes",
			args: args{
				s: "String with 'single' quotes",
			},
			want: "String with single quotes",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clearText(tt.args.s); got != tt.want {
				t.Errorf("clearText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAnimation(t *testing.T) {
	type args struct {
		config *FrameConfig
	}

	fc := &FrameConfig{
		Directives: struct {
			Type        string  "yaml:\"type\""
			Loop        int     "yaml:\"loop,omitempty\""
			Word        bool    "yaml:\"word,omitempty\""
			MaxDelay    float32 "yaml:\"maxDelay,omitempty\""
			EvenDelay   bool    "yaml:\"evenDelay,omitempty\""
			EndNewLines int     "yaml:\"endNewLines,omitempty\""
		}{
			Type: "",
		},
		Frames: []string{
			"Frame one",
		},
	}

	tests := []struct {
		name      string
		args      args
		want      IAnimation
		wantErr   bool
		animation string
	}{
		{
			name: "Return typewriter animation",
			args: args{
				config: fc,
			},
			want: &TypewriterAnimation{
				word:        false,
				maxDelay:    0,
				evenDelay:   false,
				endNewLines: 0,
				frames: []string{
					"Frame one",
				},
			},
			animation: "typewriter",
			wantErr:   false,
		},
		{
			name: "Return printer animation",
			args: args{
				config: fc,
			},
			want: &PrinterAnimation{
				maxDelay:    0,
				evenDelay:   false,
				endNewLines: 0,
				frames: []string{
					"Frame one",
				},
			},
			animation: "printer",
			wantErr:   false,
		},
		{
			name: "Return loop animation",
			args: args{
				config: fc,
			},
			want: &LoopAnimation{
				loop:        0,
				maxDelay:    0,
				endNewLines: 0,
				frames: []string{
					"Frame one",
				},
			},
			animation: "loop",
			wantErr:   false,
		},
		{
			name: "Return clear-line animation",
			args: args{
				config: fc,
			},
			want: &ClearLineAnimation{
				maxDelay:    0,
				evenDelay:   false,
				endNewLines: 0,
				frames: []string{
					"Frame one",
				},
			},
			animation: "clear-line",
			wantErr:   false,
		},
		{
			name: "Return clear-screen animation",
			args: args{
				config: fc,
			},
			want: &ClearScreenAnimation{
				maxDelay:    0,
				evenDelay:   false,
				endNewLines: 0,
				frames: []string{
					"Frame one",
				},
			},
			animation: "clear-screen",
			wantErr:   false,
		},
		{
			name: "Return error",
			args: args{
				config: fc,
			},
			want:      nil,
			animation: "notimplemented",
			wantErr:   true,
		},
	}
	for _, tt := range tests {

		tt.args.config.Directives.Type = tt.animation

		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAnimation(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAnimation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAnimation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAnimation_Render(t *testing.T) {
	tests := []struct {
		name      string
		animation IAnimation
		want      string
	}{
		{
			name: "Test ClearLineAnimation Render",
			animation: ClearLineAnimation{
				maxDelay:    0,
				evenDelay:   true,
				endNewLines: 0,
				frames: []string{
					"Frame one",
					"Frame two",
				},
			},
			want: "Frame one\rFrame two\r",
		},
		{
			name: "Test ClearScreenAnimation Render",
			animation: ClearScreenAnimation{
				maxDelay:    0,
				evenDelay:   true,
				endNewLines: 1,
				frames: []string{
					"Frame one",
					"Frame two",
				},
			},
			want: "Frame oneFrame two\n",
		},
		{
			name: "Test LoopAnimation Render",
			animation: LoopAnimation{
				loop:        3,
				maxDelay:    0,
				endNewLines: 1,
				frames: []string{
					"Frame one",
				},
			},
			want: "Frame one\rFrame one\rFrame one\r\n",
		},
		{
			name: "Test PrinterAnimation Render",
			animation: PrinterAnimation{
				maxDelay:    0,
				evenDelay:   true,
				endNewLines: 0,
				frames: []string{
					"Frame one",
					"Frame two",
					"Frame three",
				},
			},
			want: "Frame one\nFrame two\nFrame three\n",
		},
		{
			name: "Test TypewriterAnimation word Render",
			animation: TypewriterAnimation{
				word:        true,
				maxDelay:    0,
				evenDelay:   true,
				endNewLines: 0,
				frames: []string{
					"This is a typewriter word",
				},
			},
			want: "This is a typewriter word \n",
		},
		{
			name: "Test TypewriterAnimation Render",
			animation: TypewriterAnimation{
				word:        false,
				maxDelay:    0,
				evenDelay:   false,
				endNewLines: 0,
				frames: []string{
					"This is a typewriter animation",
				},
			},
			want: "This is a typewriter animation\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.animation

			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			a.Render()

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			if string(out) != tt.want {
				t.Errorf("%T Render() = %#v, want '%#v'", a, string(out), tt.want)
			}

		})
	}
}
