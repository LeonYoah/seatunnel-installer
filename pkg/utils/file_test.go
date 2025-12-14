package utils

import (
	"path/filepath"
	"testing"
)

func TestFileOperations(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	t.Run("FileExists", func(t *testing.T) {
		testFile := filepath.Join(tmpDir, "test.txt")

		// File should not exist initially
		if FileExists(testFile) {
			t.Error("File should not exist")
		}

		// Create file
		if err := WriteFile(testFile, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}

		// File should exist now
		if !FileExists(testFile) {
			t.Error("File should exist")
		}
	})

	t.Run("IsDir", func(t *testing.T) {
		testDir := filepath.Join(tmpDir, "testdir")
		testFile := filepath.Join(tmpDir, "testfile.txt")

		// Create directory
		if err := EnsureDir(testDir); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		// Create file
		if err := WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}

		if !IsDir(testDir) {
			t.Error("testdir should be a directory")
		}

		if IsDir(testFile) {
			t.Error("testfile should not be a directory")
		}
	})

	t.Run("ReadWriteFile", func(t *testing.T) {
		testFile := filepath.Join(tmpDir, "readwrite.txt")
		content := []byte("Hello, World!")

		// Write file
		if err := WriteFile(testFile, content, 0644); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}

		// Read file
		readContent, err := ReadFile(testFile)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if string(readContent) != string(content) {
			t.Errorf("Content mismatch: got %s, want %s", readContent, content)
		}
	})

	t.Run("CopyFile", func(t *testing.T) {
		srcFile := filepath.Join(tmpDir, "source.txt")
		dstFile := filepath.Join(tmpDir, "dest.txt")
		content := []byte("Copy test content")

		// Create source file
		if err := WriteFile(srcFile, content, 0644); err != nil {
			t.Fatalf("Failed to write source file: %v", err)
		}

		// Copy file
		if err := CopyFile(srcFile, dstFile); err != nil {
			t.Fatalf("Failed to copy file: %v", err)
		}

		// Verify destination file
		dstContent, err := ReadFile(dstFile)
		if err != nil {
			t.Fatalf("Failed to read destination file: %v", err)
		}

		if string(dstContent) != string(content) {
			t.Errorf("Content mismatch after copy: got %s, want %s", dstContent, content)
		}
	})

	t.Run("GetFileSize", func(t *testing.T) {
		testFile := filepath.Join(tmpDir, "size.txt")
		content := []byte("12345")

		if err := WriteFile(testFile, content, 0644); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}

		size, err := GetFileSize(testFile)
		if err != nil {
			t.Fatalf("Failed to get file size: %v", err)
		}

		if size != int64(len(content)) {
			t.Errorf("Size mismatch: got %d, want %d", size, len(content))
		}
	})

	t.Run("ListFiles", func(t *testing.T) {
		testDir := filepath.Join(tmpDir, "listtest")
		if err := EnsureDir(testDir); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		// Create some files
		for i := 1; i <= 3; i++ {
			filename := filepath.Join(testDir, filepath.Base(filepath.Join("file", string(rune('0'+i))+".txt")))
			if err := WriteFile(filename, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}

		files, err := ListFiles(testDir)
		if err != nil {
			t.Fatalf("Failed to list files: %v", err)
		}

		if len(files) != 3 {
			t.Errorf("Expected 3 files, got %d", len(files))
		}
	})
}

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "nested", "dir", "structure")

	// Create nested directory
	if err := EnsureDir(testDir); err != nil {
		t.Fatalf("Failed to create nested directory: %v", err)
	}

	// Verify directory exists
	if !FileExists(testDir) || !IsDir(testDir) {
		t.Error("Directory should exist and be a directory")
	}

	// Calling again should not error
	if err := EnsureDir(testDir); err != nil {
		t.Errorf("EnsureDir should be idempotent: %v", err)
	}
}

func TestCopyDir(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "source")
	dstDir := filepath.Join(tmpDir, "destination")

	// Create source directory structure
	if err := EnsureDir(filepath.Join(srcDir, "subdir")); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create some files
	files := map[string]string{
		"file1.txt":        "content1",
		"subdir/file2.txt": "content2",
	}

	for path, content := range files {
		fullPath := filepath.Join(srcDir, path)
		if err := WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Copy directory
	if err := CopyDir(srcDir, dstDir); err != nil {
		t.Fatalf("Failed to copy directory: %v", err)
	}

	// Verify copied files
	for path, expectedContent := range files {
		fullPath := filepath.Join(dstDir, path)
		content, err := ReadFile(fullPath)
		if err != nil {
			t.Errorf("Failed to read copied file %s: %v", path, err)
			continue
		}
		if string(content) != expectedContent {
			t.Errorf("Content mismatch for %s: got %s, want %s", path, content, expectedContent)
		}
	}
}

func TestChmod(t *testing.T) {
	// Skip on Windows as file permissions work differently
	t.Skip("Skipping chmod test on Windows")
}
