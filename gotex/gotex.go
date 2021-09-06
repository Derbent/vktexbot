package gotex

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Options contains the knobs used to change gotex's behavior.
type Options struct {
	// Command is the executable to run. It defaults to "pdflatex". Set this to
	// a full path if $PATH will not be defined in your app's environment.
	Command string
	// Runs determines how many times Command is run. This is needed for
	// documents that use refrences and packages that require multiple passes.
	// If 0, gotex will automagically attempt to determine how many runs are
	// required by parsing LaTeX log output.
	Runs int

	// Texinputs is a colon-separated list of directories containing assests
	// such as image files that are needed to compile the document. It is added
	// to $TEXINPUTS for the LaTeX process.
	Texinputs string
}

// Render takes the LaTeX document to be rendered as a string. It returns the
// resulting PDF as a []byte. If there's an error, Render will leave the
// temporary directory intact so you can check the log file to see what
// happened. The error will tell you where to find it.
func Render(document string, options Options) ([]byte, error) {
	// Set default options.
	if options.Command == "" {
		options.Command = "pdflatex"
	}

	// Create the temporary directory where LaTeX will dump its ugliness.
	var dir, err = ioutil.TempDir("", "gotex-")
	if err != nil {
		return nil, err
	}
	// Clean up the temp directory.
	defer os.RemoveAll(dir)

	// Unless a number was given, don't let automagic mode run more than this
	// many times.
	var maxRuns = 5
	if options.Runs > 0 {
		maxRuns = options.Runs
	}
	// Keep running until the document is finished or we hit an arbitrary limit.
	var runs int
	for rerun := true; rerun && runs < maxRuns; runs++ {
		err = runLatex(document, options, dir)
		if err != nil {
			return nil, err
		}
		// If in automagic mode, determine whether we need to run again.
		if options.Runs == 0 {
			rerun = needsRerun(dir)
		}
	}

	// Slurp the output.
	output, err := ioutil.ReadFile(path.Join(dir, "gotex.pdf"))
	if err != nil {
		return nil, err
	}

	return output, nil
}

// runLatex does the actual work of spawning the child and waiting for it.
func runLatex(document string, options Options, dir string) error {
	var args = []string{"-jobname=gotex", "-halt-on-error"}

	// Prepare the command.
	var cmd = exec.Command(options.Command, args...)
	// Set the cwd to the temporary directory; LaTeX will write all files there.
	cmd.Dir = dir
	// Feed the document to LaTeX over stdin.
	cmd.Stdin = strings.NewReader(document)

	// Set $TEXINPUTS if requested. The trailing colon means that LaTeX should
	// include the normal asset directories as well.
	if options.Texinputs != "" {
		cmd.Env = append(os.Environ(), "TEXINPUTS="+options.Texinputs+":")
	}

	// Launch and let it finish.
	var err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		dirToErr := path.Join(dir, "gotex.log")
		f, errOpen := os.Open(dirToErr)
		if errOpen != nil {
			return errors.New("Cannot read log file")
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "!") {
				return errors.New(line)
			}
		}
		return errors.New("Couldn't find the cause of the error")
	}
	return nil
}

// Parse the log file and attempt to determine whether another run is necessary
// to finish the document.
func needsRerun(dir string) bool {
	var file, err = os.Open(path.Join(dir, "gotex.log"))
	if err != nil {
		return false
	}
	defer file.Close()
	var scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		// Look for a line like:
		// "Label(s) may have changed. Rerun to get cross-references right."
		if strings.Contains(scanner.Text(), "Rerun to get") {
			return true
		}
	}
	return false
}
