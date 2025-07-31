package cli_test

import (
	"flag"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/mrizkifadil26/medix/utils/cli"
)

type CLIFlags struct {
	Name    string  `flag:"name" help:"Your name"`
	Age     int     `flag:"age" help:"Your age"`
	Debug   bool    `flag:"debug" help:"Enable debug mode"`
	Ignored string  // no flag tag
	SkipMe  float64 `flag:"skip" help:"unsupported type"`
}

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func TestParseCLI(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		want        CLIFlags
		expectError bool
	}{
		{
			name: "basic values",
			args: []string{"--name=Rizki", "--age=30", "--debug"},
			want: CLIFlags{Name: "Rizki", Age: 30, Debug: true},
		},
		{
			name: "default values",
			args: []string{},
			want: CLIFlags{}, // all default zero values
		},
		{
			name: "partial overrides",
			args: []string{"--debug"},
			want: CLIFlags{Debug: true},
		},
		{
			name: "unknown flag ignored by parser",
			args: []string{"--foo=bar"},
			want: CLIFlags{},
			// flag.Parse() will error, but we don’t capture it by default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags()
			os.Args = append([]string{"cmd"}, tt.args...)

			var actual CLIFlags
			cli.Parse(&actual)

			if !reflect.DeepEqual(actual, tt.want) {
				t.Errorf("expected %+v, got %+v", tt.want, actual)
			}
		})
	}
}

func TestUnsupportedTypePrintsWarning(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "--skip=3.14"}

	// Redirect stderr to capture output
	origStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	var cliFlags CLIFlags
	cli.Parse(&cliFlags)

	w.Close()
	os.Stderr = origStderr

	var out strings.Builder
	_, err := io.Copy(&out, r)
	if err != nil {
		t.Fatalf("failed to read stderr: %v", err)
	}

	if !strings.Contains(out.String(), "unsupported CLI flag type") {
		t.Errorf("expected warning for unsupported type, got: %s", out.String())
	}
}
