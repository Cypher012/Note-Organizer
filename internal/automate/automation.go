package automate

func RunAutomation() error {
	if err := UserFlow(); err != nil {
		return err
	}
	return nil
}
