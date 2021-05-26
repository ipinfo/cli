package install

import (
	"errors"
	"fmt"
)

// (un)install in bash
// basically adds/remove from .bashrc:
//
// complete -C </path/to/completion/command> -o default <command>
type bash struct {
	rc string
}

func (b bash) IsInstalled(cmd string) bool {
	completeCmd, err := BashCmd(cmd)
	if err != nil {
		return false
	}

	return lineInFile(b.rc, completeCmd)
}

func (b bash) Install(cmd string) error {
	if b.IsInstalled(cmd) {
		return fmt.Errorf("already installed in %s", b.rc)
	}

	completeCmd, err := BashCmd(cmd)
	if err != nil {
		return err
	}

	fmt.Printf("installing in %s\n", b.rc)
	return appendFile(b.rc, completeCmd)
}

func (b bash) Uninstall(cmd string) error {
	if !b.IsInstalled(cmd) {
		return fmt.Errorf("not installed in %s", b.rc)
	}

	completeCmd, err := BashCmd(cmd)
	if err != nil {
		return err
	}

	return removeFromFile(b.rc, completeCmd)
}

func BashCmd(cmd string) (string, error) {
	if binPath == "" {
		return "", errors.New("err: could not get binary path")
	}

	return fmt.Sprintf("complete -C %s -o default %s", binPath, cmd), nil
}
