package centymo

import (
	"path/filepath"
	"runtime"
)

// TemplatePatterns returns glob patterns for centymo's templates.
// Uses runtime.Caller(0) to discover centymo's package directory,
// same approach as pyeza-golang and entydad.
// Consumer apps merge these patterns with pyeza + app patterns when
// initializing the renderer.
func TemplatePatterns() []string {
	dir := packageDir()
	return []string{
		filepath.Join(dir, "templates", "plan", "*.html"),
		filepath.Join(dir, "templates", "subscription", "*.html"),
		filepath.Join(dir, "templates", "product", "*.html"),
		filepath.Join(dir, "templates", "paymentcollection", "*.html"),
		filepath.Join(dir, "templates", "inventory", "*.html"),
		filepath.Join(dir, "templates", "sales", "*.html"),
		filepath.Join(dir, "templates", "pricelist", "*.html"),
	}
}

// packageDir returns the absolute directory of this source file.
func packageDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Dir(filename)
}
