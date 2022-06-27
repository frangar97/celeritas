package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var pageData = []struct {
	name          string
	renderer      string
	template      string
	errorExpected bool
	errorMessage  string
}{
	{
		"go_page",
		"go",
		"home",
		false,
		"Error renderizando template de go",
	},
	{
		"go_page_no_template",
		"go",
		"no-file",
		true,
		"Sin error al renderizar template go cuando se espera uno",
	},
	{
		"jet_page",
		"jet",
		"home",
		false,
		"Error renderizando template de jet",
	},
	{
		"jet_page_no_template",
		"jet",
		"no-file",
		true,
		"Sin error al renderizar template jet cuando se espera uno",
	},

	{
		"invalid_render",
		"foo",
		"home",
		true,
		"Sin error renderizando con un template engine inexistente",
	},
}

func TestRender_Page(t *testing.T) {

	for _, e := range pageData {
		r, err := http.NewRequest("GET", "/some-url", nil)

		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()

		testRenderer.Renderer = e.renderer
		testRenderer.RoothPath = "./testdata"

		err = testRenderer.Page(w, r, e.template, nil, nil)

		if e.errorExpected {
			if err == nil {
				t.Errorf("%s : %s", e.name, e.errorMessage)
			}
		} else {
			if err != nil {
				t.Errorf("%s : %s: %s", e.name, e.errorMessage, err.Error())
			}
		}
	}

}

func TestRender_GoPage(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/url", nil)

	if err != nil {
		t.Error(err)
	}

	testRenderer.Renderer = "go"
	testRenderer.RoothPath = "./testdata"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Error renderizando la pagina", err)
	}
}

func TestRender_JetPage(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/url", nil)

	if err != nil {
		t.Error(err)
	}

	testRenderer.Renderer = "jet"
	testRenderer.RoothPath = "./testdata"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Error renderizando la pagina", err)
	}
}
