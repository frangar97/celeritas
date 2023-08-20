package render

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
)

type Render struct {
	Renderer   string
	RootPath   string
	Secure     bool
	Port       string
	ServerName string
	JetViews   *jet.Set
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]any
	CSRFToken       string
	Port            string
	ServerName      string
	Secure          bool
}

func (c *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data any) error {
	switch strings.ToLower(c.Renderer) {
	case "go":
		return c.goPage(w, r, view, data)
	case "jet":
		return c.jetPage(w, r, view, variables, data)
	}

	return nil
}

func (c *Render) goPage(w http.ResponseWriter, r *http.Request, view string, data any) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", c.RootPath, view))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	err = tmpl.Execute(w, td)
	if err != nil {
		return err
	}

	return nil
}

func (c *Render) jetPage(w http.ResponseWriter, r *http.Request, view string, variables, data any) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	t, err := c.JetViews.GetTemplate(fmt.Sprintf("%s.jet", view))
	if err != nil {
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		return err
	}

	return nil
}
