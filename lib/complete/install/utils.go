package install

import (
	"fmt"
	"os"
	"regexp"

	"github.com/posener/script"
)

func lineInFile(path string, line string) bool {
	return script.Cat(path).Grep(regexp.MustCompile("^"+line+"$")).Wc().Lines > 0
}

func createFile(path string, content string) error {
	return script.Echo(content).ToFile(path)
}

func appendFile(path string, content string) error {
	return script.Echo(content).AppendFile(path)
}

func removeFromFile(path string, line string) error {
	backupPath := path + ".bck"
	err := script.Cat(path).ToFile(backupPath)
	if err != nil {
		return fmt.Errorf("creating backup file: %s", err)
	}

	tmp, err := script.Cat(path).Modify(script.Grep{Re: regexp.MustCompile("^" + line + "$"), Inverse: true}).ToTempFile()
	if err != nil {
		return fmt.Errorf("failed remove: %s", err)
	}
	defer os.Remove(tmp)

	err = script.Cat(tmp).ToFile(path)
	if err != nil {
		restoreErr := script.Cat(backupPath).ToFile(path)
		if restoreErr != nil {
			return fmt.Errorf("failed write: %s, and failed restore: %s", err, restoreErr)
		}
	}
	return nil
}
