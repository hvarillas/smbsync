package banner

import (
	"bytes"
	"strings"
	"testing"
)

func TestText(t *testing.T) {
	banner := Text()
	
	if banner == "" {
		t.Error("Banner text should not be empty")
	}
	
	// Check if banner contains expected elements
	expectedElements := []string{
		"v1.0",
		"Hernán Varillas",
		"hernan.varillas93@gmail.com",
	}
	
	for _, element := range expectedElements {
		if !strings.Contains(banner, element) {
			t.Errorf("Banner should contain '%s'", element)
		}
	}
}

func TestPrint(t *testing.T) {
	var buf bytes.Buffer
	Print(&buf)
	
	output := buf.String()
	if output == "" {
		t.Error("Print should produce output")
	}
	
	// Should end with newline
	if !strings.HasSuffix(output, "\n") {
		t.Error("Print output should end with newline")
	}
	
	// Should contain the same content as Text()
	expectedText := Text() + "\n"
	if output != expectedText {
		t.Error("Print output should match Text() + newline")
	}
}

func TestBannerFormat(t *testing.T) {
	banner := Text()
	lines := strings.Split(banner, "\n")
	
	// Banner should have multiple lines
	if len(lines) < 5 {
		t.Error("Banner should have multiple lines")
	}
	
	// Check for ASCII art characters
	hasAsciiArt := false
	for _, line := range lines {
		if strings.Contains(line, "█") || strings.Contains(line, "╗") || strings.Contains(line, "╚") {
			hasAsciiArt = true
			break
		}
	}
	
	if !hasAsciiArt {
		t.Error("Banner should contain ASCII art characters")
	}
}
