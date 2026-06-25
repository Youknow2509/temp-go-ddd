package initialize

// ===
// Initialize the system
func Initialize() error {
	// Initialize configurations system
	if err := initializeConfig(); err != nil {
		return err
	}
	// Initialize the logger
	if err := initializeLogger(); err != nil {
		return err
	}

	return nil
}
