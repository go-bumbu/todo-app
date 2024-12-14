package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

// SimpleText Is a simple handler that will print a list of navigation destination based on the map passed upon creation.
type SimpleText struct {
	Text  string
	Links []Link
}
type Link struct {
	Text  string
	Url   string
	Child []Link
}

func (h SimpleText) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", " text/html")
	if r.Method == http.MethodGet {

		var s strings.Builder

		s.WriteString("GET: " + h.Text)

		if len(h.Links) > 0 {
			writeRec(&s, h.Links)
		}
		fmt.Fprint(w, s.String())
		return
	}

	if r.Method == http.MethodPost {
		fmt.Fprintf(w, "POST: %s", h.Text)
		return
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func writeRec(s *strings.Builder, links []Link) {
	s.WriteString("<ul>")
	for _, link := range links {
		if link.Url == "" {
			s.WriteString(fmt.Sprintf("<li>%s", link.Text))
			if len(link.Child) > 0 {
				writeRec(s, link.Child)
			}
			s.WriteString("</li>")
		} else {
			s.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s</a>", link.Url, link.Text))
			if len(link.Child) > 0 {
				writeRec(s, link.Child)
			}
			s.WriteString("</li>")
		}
	}
	s.WriteString("</ul>")
}
