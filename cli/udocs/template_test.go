package udocs

import "testing"

func TestMustParseTemplate(t *testing.T) {
	files := DefaultTemplateFiles()
	params := map[string]interface{}{"email": "user@email.com"}

	template := MustParseTemplate(params, files...)
	if template == nil {
		t.Error("(*udocs.Template) cannot be nil")
	}

	if template.tmpl == nil {
		t.Error("(*udocs.Template).tmpl cannot be nil")
	}

	if template.files == nil {
		t.Error("(*udocs.Template).files cannot be nil")
	}
	for i := range template.files {
		if template.files[i] != files[i] {
			t.Errorf("(*udocs.Template).files %v does not match default template files", template.files)
		}
	}

	if template.params == nil {
		t.Error("(*udocs.Template).params cannot be nil")
	}
	for k := range template.params {
		if template.params[k] != params[k] {
			t.Errorf("(*udocs.Template).params %v does not match expected params", template.params)
		}
	}
}
