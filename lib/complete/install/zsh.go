package install

import (
	"errors"
	"fmt"
)

// (un)install in zsh
// basically adds/remove from .zshrc:
//
// autoload -U +X bashcompinit && bashcompinit"
// complete -o default -C </path/to/completion/command> <command>
type zsh struct {
	rc string
}

func (z zsh) IsInstalled(cmd string) bool {
	completeCmd, err := ZshCmd(cmd)
	if err != nil {
		return false
	}

	return lineInFile(z.rc, completeCmd)
}

func (z zsh) Install(cmd string) error {
	if z.IsInstalled(cmd) {
		return fmt.Errorf("already installed in %s", z.rc)
	}

	completeCmd, err := ZshCmd(cmd)
	if err != nil {
		return err
	}

	fmt.Printf("installing in %s\n", z.rc)
	return appendFile(z.rc, completeCmd)
}

func (z zsh) Uninstall(cmd string) error {
	if !z.IsInstalled(cmd) {
		return fmt.Errorf("not installed in %s", z.rc)
	}

	completeCmd, err := ZshCmd(cmd)
	if err != nil {
		return err
	}

	return removeFromFile(z.rc, completeCmd)
}

func ZshCmd(cmd string) (string, error) {
	if binPath == "" {
		return "", errors.New("err: could not get binary path")
	}

	return "autoload -U +X bashcompinit && bashcompinit" +
		"\n" + fmt.Sprintf("complete -o default -C %s %s", binPath, cmd), nil
}
