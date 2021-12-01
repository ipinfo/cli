package main

func saveToken(tok string) error {
	// write token.
	gConfig.Token = tok
	err := SetConfig(gConfig)
	if err != nil {
		return err
	}

	return nil
}

func deleteToken() error {
	gConfig.Token = ""
	err := SetConfig(gConfig)
	if err != nil {
		return err
	}
	return nil
}

func restoreToken() string {
	return string(gConfig.Token)
}
