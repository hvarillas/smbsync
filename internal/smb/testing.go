package smb

// GetRegexFiles exposes the internal getRegexFiles function for testing
func GetRegexFiles(regex, path string) []string {
	return getRegexFiles(regex, path)
}
