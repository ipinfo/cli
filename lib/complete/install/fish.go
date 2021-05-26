package install

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// (un)install in fish

type fish struct {
	configDir string
}

func (f fish) IsInstalled(cmd string) bool {
	completionFile := f.getCompletionFilePath(cmd)
	if _, err := os.Stat(completionFile); err == nil {
		return true
	}
	return false
}

func (f fish) Install(cmd string) error {
	if f.IsInstalled(cmd) {
		return fmt.Errorf(
			"already installed at %s",
			f.getCompletionFilePath(cmd),
		)
	}

	completionFile := f.getCompletionFilePath(cmd)
	completeCmd, err := FishCmd(cmd)
	if err != nil {
		return err
	}

	fmt.Printf("installing in %s\n", completionFile)
	return createFile(completionFile, completeCmd)
}

func (f fish) Uninstall(cmd string) error {
	if !f.IsInstalled(cmd) {
		return fmt.Errorf("does not installed in %s", f.configDir)
	}

	completionFile := f.getCompletionFilePath(cmd)
	return os.Remove(completionFile)
}

func (f fish) getCompletionFilePath(cmd string) string {
	return filepath.Join(
		f.configDir,
		"completions",
		fmt.Sprintf("%s.fish", cmd),
	)
}

func FishCmd(cmd string) (string, error) {
	var buf bytes.Buffer

	if binPath == "" {
		return "", errors.New("err: could not get binary path")
	}

	params := struct{ Cmd, Bin string }{cmd, binPath}
	tmpl := template.Must(template.New("cmd").Parse(`function __complete_{{.Cmd}}
    set -lx COMP_LINE (commandline -cp)
    test -z (commandline -ct)
    and set COMP_LINE "$COMP_LINE "
    {{.Bin}}
end
complete -f -c {{.Cmd}} -a "(__complete_{{.Cmd}})"`))
	err := tmpl.Execute(&buf, params)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
