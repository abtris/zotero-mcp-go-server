package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/abtris/zotero-go-client/zotero"
)

func testItem() *zotero.Item {
	return &zotero.Item{
		Key:     "CIT1",
		Version: 3,
		Data: zotero.ItemData{
			ItemType: "journalArticle",
			Title:    "Deep Learning for NLP",
			Creators: []zotero.Creator{
				{CreatorType: "author", FirstName: "Jane", LastName: "Doe"},
				{CreatorType: "author", FirstName: "John", LastName: "Smith"},
			},
			Date: "2024",
		},
	}
}

func TestFormatCitationAPA(t *testing.T) {
	item := testItem()
	got := formatCitationAPA(item)
	if !strings.Contains(got, "Doe, Jane") {
		t.Errorf("APA should contain 'Doe, Jane', got: %s", got)
	}
	if !strings.Contains(got, "(2024)") {
		t.Errorf("APA should contain '(2024)', got: %s", got)
	}
	if !strings.Contains(got, "Deep Learning for NLP") {
		t.Errorf("APA should contain title, got: %s", got)
	}
}

func TestFormatCitationChicago(t *testing.T) {
	got := formatCitationChicago(testItem())
	if !strings.Contains(got, "Jane Doe") {
		t.Errorf("Chicago should contain 'Jane Doe', got: %s", got)
	}
	if !strings.Contains(got, "2024") {
		t.Errorf("Chicago should contain year, got: %s", got)
	}
}

func TestFormatCitationMLA(t *testing.T) {
	got := formatCitationMLA(testItem())
	if !strings.Contains(got, "Doe, Jane") {
		t.Errorf("MLA should contain 'Doe, Jane', got: %s", got)
	}
	if !strings.Contains(got, "\"Deep Learning for NLP.\"") {
		t.Errorf("MLA should quote the title, got: %s", got)
	}
}

func TestFormatCitationHarvard(t *testing.T) {
	got := formatCitationHarvard(testItem())
	if !strings.Contains(got, "Doe") {
		t.Errorf("Harvard should contain last name, got: %s", got)
	}
	if !strings.Contains(got, "(2024)") {
		t.Errorf("Harvard should contain '(2024)', got: %s", got)
	}
}

func TestFormatCitationIEEE(t *testing.T) {
	got := formatCitationIEEE(testItem())
	if !strings.Contains(got, "Jane Doe") {
		t.Errorf("IEEE should contain 'Jane Doe', got: %s", got)
	}
	if !strings.Contains(got, "\"Deep Learning for NLP,\"") {
		t.Errorf("IEEE should quote the title, got: %s", got)
	}
}

func TestFormatCitationBibTeX(t *testing.T) {
	got := formatCitationBibTeX(testItem())
	if !strings.Contains(got, "@article{CIT1") {
		t.Errorf("BibTeX should start with @article{CIT1, got: %s", got)
	}
	if !strings.Contains(got, "author = {Jane Doe") {
		t.Errorf("BibTeX should contain author, got: %s", got)
	}
	if !strings.Contains(got, "year = {2024}") {
		t.Errorf("BibTeX should contain year, got: %s", got)
	}
}

func TestFormatCitationNoDate(t *testing.T) {
	item := testItem()
	item.Data.Date = ""
	got := formatCitationAPA(item)
	if !strings.Contains(got, "n.d.") {
		t.Errorf("APA with no date should contain 'n.d.', got: %s", got)
	}
}

func TestFormatCitationDispatch(t *testing.T) {
	item := testItem()
	tests := []struct {
		style    string
		contains string
	}{
		{"apa", "(2024)"},
		{"chicago", "Jane Doe"},
		{"mla", "\"Deep Learning"},
		{"harvard", "Doe, Smith"},
		{"ieee", "Jane Doe"},
		{"bibtex", "@article"},
		{"unknown", "(2024)"}, // defaults to APA
	}
	for _, tt := range tests {
		got := formatCitation(item, tt.style)
		if !strings.Contains(got, tt.contains) {
			t.Errorf("formatCitation(%q) should contain %q, got: %s", tt.style, tt.contains, got)
		}
	}
}

func TestGenerateCitationToolAPI(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/CIT1" {
			http.NotFound(w, r)
			return
		}
		q := r.URL.Query()
		if q.Get("format") == "bib" {
			if q.Get("style") != "apa" {
				t.Errorf("expected style=apa, got %q", q.Get("style"))
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<div class="csl-entry">Doe, J., &amp; Smith, J. (2024). Deep Learning for NLP.</div>`))
			return
		}
		json.NewEncoder(w).Encode(testItem())
	})

	text := callTool(t, session, "generate_citation", map[string]any{
		"item_key": "CIT1",
		"style":    "apa",
	})
	if !strings.Contains(text, "Doe") {
		t.Errorf("expected citation with author, got: %s", text)
	}
	if !strings.Contains(text, "2024") {
		t.Errorf("expected citation with year, got: %s", text)
	}
}

func TestGenerateCitationToolBibTeX(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/CIT1" {
			http.NotFound(w, r)
			return
		}
		// BibTeX uses Get (JSON), not GetBibliography
		json.NewEncoder(w).Encode(testItem())
	})

	text := callTool(t, session, "generate_citation", map[string]any{
		"item_key": "CIT1",
		"style":    "bibtex",
	})
	if !strings.Contains(text, "@article{CIT1") {
		t.Errorf("expected BibTeX entry, got: %s", text)
	}
}

func TestGenerateCitationToolWithLocale(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/CIT1" {
			http.NotFound(w, r)
			return
		}
		q := r.URL.Query()
		if q.Get("format") == "bib" {
			if q.Get("locale") != "de-DE" {
				t.Errorf("expected locale=de-DE, got %q", q.Get("locale"))
			}
			w.Write([]byte(`<div class="csl-entry">Doe &amp; Smith (2024). Deep Learning for NLP.</div>`))
			return
		}
		json.NewEncoder(w).Encode(testItem())
	})

	text := callTool(t, session, "generate_citation", map[string]any{
		"item_key": "CIT1",
		"style":    "apa",
		"locale":   "de-DE",
	})
	if !strings.Contains(text, "Doe") {
		t.Errorf("expected citation with author, got: %s", text)
	}
}

func TestGenerateCitationToolCSLStyleMapping(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/CIT1" {
			http.NotFound(w, r)
			return
		}
		q := r.URL.Query()
		if q.Get("format") == "bib" {
			if q.Get("style") != "chicago-note-bibliography" {
				t.Errorf("expected chicago CSL ID, got %q", q.Get("style"))
			}
			w.Write([]byte(`<div class="csl-entry">Chicago citation here.</div>`))
			return
		}
		json.NewEncoder(w).Encode(testItem())
	})

	text := callTool(t, session, "generate_citation", map[string]any{
		"item_key": "CIT1",
		"style":    "chicago",
	})
	if !strings.Contains(text, "Chicago citation") {
		t.Errorf("expected Chicago citation, got: %s", text)
	}
}
