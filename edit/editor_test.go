package edit

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestEditor_ReplaceSection_CreateNewFile(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")

	// Create editor
	editor := NewEditor(filePath)

	// Replace section (file doesn't exist yet)
	content := []byte("# Generated Content\n\nThis is the generated documentation.\n")
	if err := editor.ReplaceSection(content); err != nil {
		t.Fatalf("ReplaceSection failed: %v", err)
	}

	// Read result
	result, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	// Verify structure
	resultStr := string(result)
	if !strings.Contains(resultStr, MarkerBegin) {
		t.Errorf("Result missing begin marker")
	}
	if !strings.Contains(resultStr, MarkerEnd) {
		t.Errorf("Result missing end marker")
	}
	if !strings.Contains(resultStr, "# Generated Content") {
		t.Errorf("Result missing generated content")
	}

	// Verify order
	beginIdx := strings.Index(resultStr, MarkerBegin)
	endIdx := strings.Index(resultStr, MarkerEnd)
	contentIdx := strings.Index(resultStr, "# Generated Content")
	if beginIdx >= contentIdx || contentIdx >= endIdx {
		t.Errorf("Markers not in correct order: begin=%d, content=%d, end=%d", beginIdx, contentIdx, endIdx)
	}
}

func TestEditor_ReplaceSection_BasicReplacement(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")

	// Create initial file with markers
	initialContent := `# My README

Some intro text.

<!--envdoc:begin-->
Old generated content
<!--envdoc:end-->

Some footer text.
`
	if err := os.WriteFile(filePath, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	// Create editor
	editor := NewEditor(filePath)

	// Replace section
	newContent := []byte("# New Generated Content\n\nThis is updated.\n")
	if err := editor.ReplaceSection(newContent); err != nil {
		t.Fatalf("ReplaceSection failed: %v", err)
	}

	// Read result
	result, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	resultStr := string(result)

	// Verify content before markers is preserved
	if !strings.Contains(resultStr, "# My README") {
		t.Errorf("Content before markers was not preserved")
	}
	if !strings.Contains(resultStr, "Some intro text.") {
		t.Errorf("Content before markers was not preserved")
	}

	// Verify new content is present
	if !strings.Contains(resultStr, "# New Generated Content") {
		t.Errorf("New content not found")
	}
	if !strings.Contains(resultStr, "This is updated.") {
		t.Errorf("New content not found")
	}

	// Verify old content is gone
	if strings.Contains(resultStr, "Old generated content") {
		t.Errorf("Old content was not replaced")
	}

	// Verify content after markers is preserved
	if !strings.Contains(resultStr, "Some footer text.") {
		t.Errorf("Content after markers was not preserved")
	}
}

func TestEditor_ReplaceSection_EmptyContent(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")

	// Create initial file with markers
	initialContent := `# README

<!--envdoc:begin-->
Content to be removed
<!--envdoc:end-->

Footer.
`
	if err := os.WriteFile(filePath, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	// Create editor
	editor := NewEditor(filePath)

	// Replace with empty content
	if err := editor.ReplaceSection([]byte{}); err != nil {
		t.Fatalf("ReplaceSection failed: %v", err)
	}

	// Read result
	result, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	resultStr := string(result)

	// Verify markers still present
	if !strings.Contains(resultStr, MarkerBegin) {
		t.Errorf("Begin marker missing")
	}
	if !strings.Contains(resultStr, MarkerEnd) {
		t.Errorf("End marker missing")
	}

	// Verify old content is gone
	if strings.Contains(resultStr, "Content to be removed") {
		t.Errorf("Old content was not removed")
	}

	// Verify other content preserved
	if !strings.Contains(resultStr, "# README") {
		t.Errorf("Header not preserved")
	}
	if !strings.Contains(resultStr, "Footer.") {
		t.Errorf("Footer not preserved")
	}
}

func TestEditor_ReplaceSection_MissingMarkers(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{
			name:    "no markers",
			content: "# README\n\nNo markers here.\n",
			wantErr: "no markers found",
		},
		{
			name:    "only begin marker",
			content: "# README\n\n<!--envdoc:begin-->\n\nNo end marker.\n",
			wantErr: "missing end marker",
		},
		{
			name:    "only end marker",
			content: "# README\n\nNo begin marker.\n\n<!--envdoc:end-->\n",
			wantErr: "missing begin marker",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "test.md")

			// Create file with test content
			if err := os.WriteFile(filePath, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Create editor
			editor := NewEditor(filePath)

			// Attempt to replace section
			err := editor.ReplaceSection([]byte("new content"))
			if err == nil {
				t.Fatalf("Expected error but got none")
			}

			// Verify error message
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Expected error containing %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestEditor_ReplaceSection_MarkersWrongOrder(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")

	// Create file with markers in wrong order
	content := `# README

<!--envdoc:end-->
Content
<!--envdoc:begin-->

Footer.
`
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create editor
	editor := NewEditor(filePath)

	// Attempt to replace section
	err := editor.ReplaceSection([]byte("new content"))
	if err == nil {
		t.Fatalf("Expected error but got none")
	}

	// Verify error message
	if !strings.Contains(err.Error(), "wrong order") {
		t.Errorf("Expected error about wrong order, got: %v", err)
	}
}

func TestEditor_ReplaceSection_DuplicateMarkers(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{
			name: "duplicate begin markers",
			content: `# README

<!--envdoc:begin-->
Section 1
<!--envdoc:end-->

<!--envdoc:begin-->
Section 2
<!--envdoc:end-->
`,
			wantErr: "duplicate begin markers",
		},
		{
			name: "duplicate end markers",
			content: `# README

<!--envdoc:begin-->
Section
<!--envdoc:end-->
More content
<!--envdoc:end-->
`,
			wantErr: "duplicate end markers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "test.md")

			// Create file with test content
			if err := os.WriteFile(filePath, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Create editor
			editor := NewEditor(filePath)

			// Attempt to replace section
			err := editor.ReplaceSection([]byte("new content"))
			if err == nil {
				t.Fatalf("Expected error but got none")
			}

			// Verify error message
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Expected error containing %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestEditor_ReplaceSection_MultilineContent(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")

	// Create initial file
	initialContent := `# Project README

This is the introduction.

<!--envdoc:begin-->
Old docs
<!--envdoc:end-->

## Other Section

This should be preserved.
`
	if err := os.WriteFile(filePath, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	// Create editor
	editor := NewEditor(filePath)

	// Replace with multiline content
	newContent := `## Environment Variables

### PORT

**Type**: int

**Default**: 8080

**Description**: Server port

### DATABASE_URL

**Type**: string

**Required**: Yes

**Description**: Database connection string
`
	if err := editor.ReplaceSection([]byte(newContent)); err != nil {
		t.Fatalf("ReplaceSection failed: %v", err)
	}

	// Read result
	result, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	resultStr := string(result)

	// Verify all sections present
	requiredSections := []string{
		"# Project README",
		"This is the introduction.",
		MarkerBegin,
		"## Environment Variables",
		"### PORT",
		"### DATABASE_URL",
		MarkerEnd,
		"## Other Section",
		"This should be preserved.",
	}

	for _, section := range requiredSections {
		if !strings.Contains(resultStr, section) {
			t.Errorf("Missing required section: %q", section)
		}
	}

	// Verify old content is gone
	if strings.Contains(resultStr, "Old docs") {
		t.Errorf("Old content was not replaced")
	}
}

func TestEditor_ReplaceSection_CreateInSubdirectory(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "subdir", "nested", "test.md")

	// Create editor
	editor := NewEditor(filePath)

	// Replace section (directories don't exist yet)
	content := []byte("# Content\n")
	if err := editor.ReplaceSection(content); err != nil {
		t.Fatalf("ReplaceSection failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("File was not created: %v", err)
	}

	// Read and verify content
	result, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	resultStr := string(result)
	if !strings.Contains(resultStr, MarkerBegin) || !strings.Contains(resultStr, MarkerEnd) {
		t.Errorf("Markers not found in created file")
	}
	if !strings.Contains(resultStr, "# Content") {
		t.Errorf("Content not found in created file")
	}
}

func TestEditor_ReplaceSection_ContentWithoutTrailingNewline(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")

	// Create initial file with markers
	initialContent := `# README

<!--envdoc:begin-->
Old content
<!--envdoc:end-->

Footer.
`
	if err := os.WriteFile(filePath, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	// Create editor
	editor := NewEditor(filePath)

	// Replace with content that doesn't end with newline
	newContent := []byte("New content without trailing newline")
	if err := editor.ReplaceSection(newContent); err != nil {
		t.Fatalf("ReplaceSection failed: %v", err)
	}

	// Read result
	result, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	resultStr := string(result)

	// Verify new content is present
	if !strings.Contains(resultStr, "New content without trailing newline") {
		t.Errorf("New content not found")
	}

	// Verify markers still present
	if !strings.Contains(resultStr, MarkerBegin) {
		t.Errorf("Begin marker missing")
	}
	if !strings.Contains(resultStr, MarkerEnd) {
		t.Errorf("End marker missing")
	}

	// Verify footer is preserved
	if !strings.Contains(resultStr, "Footer.") {
		t.Errorf("Footer not preserved")
	}
}

func TestEditor_ReplaceSection_CreateFileWithoutTrailingNewline(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")

	// Create editor
	editor := NewEditor(filePath)

	// Create file with content that doesn't end with newline
	content := []byte("Content without newline")
	if err := editor.ReplaceSection(content); err != nil {
		t.Fatalf("ReplaceSection failed: %v", err)
	}

	// Read result
	result, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	resultStr := string(result)

	// Verify content and markers are present
	if !strings.Contains(resultStr, "Content without newline") {
		t.Errorf("Content not found")
	}
	if !strings.Contains(resultStr, MarkerBegin) {
		t.Errorf("Begin marker missing")
	}
	if !strings.Contains(resultStr, MarkerEnd) {
		t.Errorf("End marker missing")
	}

	// Verify proper structure (newline added after content)
	if !strings.Contains(resultStr, "Content without newline\n"+MarkerEnd) {
		t.Errorf("Newline was not added after content")
	}
}

func TestEditor_ReplaceSection_StatError(t *testing.T) {
	// This test is tricky to implement portably since we need to trigger
	// a stat error that's not os.IsNotExist. We can use a file in a directory
	// without permissions (on Unix systems)
	if runtime.GOOS == "windows" {
		t.Skip("Unix file permissions are not enforced on Windows")
	}
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	tmpDir := t.TempDir()

	// Create a directory without read/execute permissions
	restrictedDir := filepath.Join(tmpDir, "restricted")
	if err := os.Mkdir(restrictedDir, 0o000); err != nil {
		t.Fatalf("Failed to create restricted directory: %v", err)
	}
	defer func() {
		_ = os.Chmod(restrictedDir, 0o755) // Restore permissions for cleanup
	}()

	filePath := filepath.Join(restrictedDir, "test.md")
	editor := NewEditor(filePath)

	// Attempt to replace section should fail with stat error
	err := editor.ReplaceSection([]byte("content"))
	if err == nil {
		t.Fatalf("Expected error but got none")
	}

	if !strings.Contains(err.Error(), "failed to stat file") {
		t.Errorf("Expected stat error, got: %v", err)
	}
}

func TestEditor_ReplaceSection_MkdirError(t *testing.T) {
	// Create a file where we want to create a directory
	// This will cause an error when trying to create nested directories
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	tmpDir := t.TempDir()

	// Create a regular file
	blockingFile := filepath.Join(tmpDir, "file")
	if err := os.WriteFile(blockingFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create blocking file: %v", err)
	}

	// Try to create a file where we need "file" to be a directory
	filePath := filepath.Join(blockingFile, "test.md")
	editor := NewEditor(filePath)

	// Attempt to replace section should fail
	err := editor.ReplaceSection([]byte("content"))
	if err == nil {
		t.Fatalf("Expected error but got none")
	}

	// The error can be either a stat error (not a directory) or directory creation error
	// Both are valid since they both prevent creating the file
	errMsg := err.Error()
	hasValidError := strings.Contains(errMsg, "failed to create directory") ||
		strings.Contains(errMsg, "failed to create temp file") ||
		strings.Contains(errMsg, "failed to stat file") ||
		strings.Contains(errMsg, "not a directory")

	if !hasValidError {
		t.Errorf("Expected file system error, got: %v", err)
	}
}

func TestEditor_ReplaceSection_WriteToReadOnlyDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix file permissions are not enforced on Windows")
	}
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0o755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Create initial file
	filePath := filepath.Join(readOnlyDir, "test.md")
	initialContent := `# README

<!--envdoc:begin-->
Old content
<!--envdoc:end-->

Footer.
`
	if err := os.WriteFile(filePath, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	// Make directory read-only
	if err := os.Chmod(readOnlyDir, 0o555); err != nil {
		t.Fatalf("Failed to change directory permissions: %v", err)
	}
	defer func() {
		_ = os.Chmod(readOnlyDir, 0o755) // Restore for cleanup
	}()

	editor := NewEditor(filePath)

	// Attempt to replace section should fail (can't create temp file)
	err := editor.ReplaceSection([]byte("new content"))
	if err == nil {
		t.Fatalf("Expected error but got none")
	}

	if !strings.Contains(err.Error(), "failed to create temp file") {
		t.Errorf("Expected temp file creation error, got: %v", err)
	}
}
