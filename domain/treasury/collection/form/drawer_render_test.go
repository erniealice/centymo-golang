package form_test

import (
	"net/http/httptest"
	"testing"
	"testing/fstest"

	pyeza "github.com/erniealice/pyeza-golang"

	"github.com/erniealice/centymo-golang/domain/treasury/collection"
	"github.com/erniealice/centymo-golang/domain/treasury/collection/form"
)

// The Add/Edit Collection drawer binds form.Data.Labels (collection.FormLabels) directly; a
// `.Labels.Form.X` path (left by the 20260611 package split) failed every render with
// "can't evaluate field Form in type collection.FormLabels" (found in the S1 E2E wave).
func TestCollectionDrawerFormRenders(t *testing.T) {
	shell := fstest.MapFS{"app-shell.html": {Data: []byte(`{{define "app-shell"}}[shell]{{end}}`)}}
	r := pyeza.NewHTMLRendererFromFS(pyeza.SharedFS, shell, collection.TemplatesFS)
	if err := r.Init(); err != nil {
		t.Fatalf("template set does not parse: %v", err)
	}
	w := httptest.NewRecorder()
	data := &form.Data{Labels: collection.DefaultLabels().Form, CommonLabels: pyeza.CommonLabels{}}
	if err := r.Render(w, "collection-drawer-form", data); err != nil {
		t.Fatalf("render collection-drawer-form: %v", err)
	}
	if w.Body.Len() == 0 {
		t.Fatal("empty render")
	}
}
