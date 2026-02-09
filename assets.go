package centymo

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// CopyStyles copies centymo's CSS assets to the target directory.
// Uses runtime.Caller(0) via packageDir() to discover centymo's package
// directory, same approach as pyeza-golang.
//
// Files are copied to {targetDir}/centymo/ to keep them namespaced.
//
// Example:
//
//	cssDir := filepath.Join("assets", "css")
//	if err := centymo.CopyStyles(cssDir); err != nil {
//	    log.Printf("Warning: Failed to copy centymo styles: %v", err)
//	}
func CopyStyles(targetDir string) error {
	dir := packageDir()
	if dir == "" {
		return fmt.Errorf("could not determine centymo package directory")
	}

	srcDir := filepath.Join(dir, "assets", "css")
	dstDir := filepath.Join(targetDir, "centymo")

	copied, err := copyDirFiles(srcDir, dstDir, "*.css")
	if err != nil {
		return fmt.Errorf("failed to copy centymo styles: %w", err)
	}

	if copied == 0 {
		log.Printf("centymo: no CSS files found in %s", srcDir)
		return nil
	}

	log.Printf("Copied %d centymo styles to: %s", copied, dstDir)
	return nil
}

// CopyStaticAssets copies centymo's JavaScript assets to the target directory.
// Uses runtime.Caller(0) via packageDir() to discover centymo's package
// directory, same approach as pyeza-golang.
//
// Files are copied to {targetDir}/centymo/ to keep them namespaced.
//
// Example:
//
//	jsDir := filepath.Join("assets", "js")
//	if err := centymo.CopyStaticAssets(jsDir); err != nil {
//	    log.Printf("Warning: Failed to copy centymo assets: %v", err)
//	}
func CopyStaticAssets(targetDir string) error {
	dir := packageDir()
	if dir == "" {
		return fmt.Errorf("could not determine centymo package directory")
	}

	srcDir := filepath.Join(dir, "assets", "js")
	dstDir := filepath.Join(targetDir, "centymo")

	copied, err := copyDirFiles(srcDir, dstDir, "*.js")
	if err != nil {
		return fmt.Errorf("failed to copy centymo assets: %w", err)
	}

	if copied == 0 {
		log.Printf("centymo: no JS files found in %s", srcDir)
		return nil
	}

	log.Printf("Copied %d centymo assets to: %s", copied, dstDir)
	return nil
}

// copyDirFiles copies all files matching a glob pattern from srcDir to dstDir.
func copyDirFiles(srcDir, dstDir, pattern string) (int, error) {
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create target directory: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(srcDir, pattern))
	if err != nil {
		return 0, fmt.Errorf("failed to list source files: %w", err)
	}

	var copied int
	for _, srcFile := range files {
		data, err := os.ReadFile(srcFile)
		if err != nil {
			log.Printf("Warning: Failed to read %s: %v", srcFile, err)
			continue
		}

		dstFile := filepath.Join(dstDir, filepath.Base(srcFile))
		if err := os.WriteFile(dstFile, data, 0644); err != nil {
			return copied, fmt.Errorf("failed to write %s: %w", dstFile, err)
		}
		copied++
	}

	return copied, nil
}
