package lib

func CmdIP2n(cmd string) (string, error) {
	res, err := IPtoDecimalStr(cmd)
	if err != nil {
		return "", err
	}

	return res, err
}
