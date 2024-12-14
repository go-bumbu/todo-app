package handlers

import (
	"html/template"
	"net/http"
)

// TODO, possible improvements
// * template does not need to be parsed on every request
// * add additional template functions, e.g. https://gist.github.com/alex-leonhardt/8ed3f78545706d89d466434fb6870023
// * if need be extend into own templating library for html data, e.g. convention folders and rendering cycle
//   probably does not make much sense with an SPA

const tmpName = "banana"

func TemplateWithRequest(inTmpl string, data func(r *http.Request) map[string]interface{}) func(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.New(tmpName).Parse(inTmpl)
	return func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{}
		if data != nil {
			payload = data(r)
		}
		err := tmpl.Execute(w, &payload)
		if err != nil {
			http.Error(w, "templating error: ", http.StatusInternalServerError)
		}
	}
}
