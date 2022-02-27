package plotocrat

import (
	"bytes"
	"crypto/rand"
	"database/sql"

	// To make go:embed work
	_ "embed"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/ConradIrwin/plotocrat/env"
)

// Chart stored the data for a chart, just an axis label and a data blob that was uploaded.
type Chart struct {
	ID string `db:"id"`
	XAxis string `db:"x_axis"`
	Data string `db:"data"`
}

func readSinglePart(r *http.Request) (string, string, error) {
	reader, err := r.MultipartReader()
	if err != nil {
		return "", "", err
	}

	part, err := reader.NextPart()
	if err != nil {
		return "", "", err
	}

	buf, err := io.ReadAll(io.LimitReader(part, 10 * 1024 * 1024 + 1))
	if err != nil {
		return "", "", err
	}
	if len(buf) == 10 * 1024 * 1024 + 1 {
		return "", "", fmt.Errorf("file is more than 10Mb")
	}
	
	if _, err = reader.NextPart(); err != io.EOF {
		return "", "", fmt.Errorf("too many parts")
	}

	if !utf8.ValidString(part.FormName()) {
		return "", "", fmt.Errorf("filename is not utf8")
	}

	if !utf8.Valid(buf) {
		return "", "", fmt.Errorf("file content is not utf8")
	}

	return part.FormName(), string(buf), nil
}

func newID() string {
	buf := make([]byte,10)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}

func isValidID(id string) bool {
	b, err := hex.DecodeString(id)
	return err == nil && len(b) == 10
}

// Handler handles all requests to plotocrat.com
// TODO: using a framework would be nice...
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/" {
		name, data, err := readSinglePart(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad-request: " + err.Error() + "\n\n"))
			w.Write([]byte("Usage: curl -F \"x-axis-title=@filename\" plotocrat.com\n"))
			return
		}

		chart := &Chart{
			ID: newID(),
			XAxis: name,
			Data: data,
		}

		if _, err := env.DB.NamedExec("INSERT INTO charts (id, x_axis, data) VALUES (:id, :x_axis, :data)", chart); err != nil {
			panic(err)
		}

		fmt.Fprintln(w, env.Config.ServerOrigin + "/" + chart.ID)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad-request: GET is only supported method\n\n"))
		w.Write([]byte("Usage: curl -F \"x-axis-title=@filename\" plotocrat.com\n"))
		return 
	}

	if r.URL.Path == "/" {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Write(getIndexRendered())
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	if !isValidID(path) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	chart := Chart{}
	err := env.DB.Get(&chart, "SELECT * FROM charts WHERE id = $1;", path)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	renderChart(w, &chart)
}

//go:embed chart.html.tmpl
var chartHTMLTemplate string
var chartHTMLCompiled *template.Template

//go:embed index.html.tmpl
var indexHTMLTemplate string
var indexHTMLRendered []byte

//go:embed index.css
var indexCSS template.CSS
//go:embed index.js
var indexJS template.JS
var once sync.Once

func compile() {
	work := func () {
		t, err := template.New("chart").Parse(chartHTMLTemplate)
		if err != nil {
			panic(err)
		}
		chartHTMLCompiled = t

		t, err = template.New("chart").Parse(indexHTMLTemplate)
		if err != nil {
			panic(err)
		}
		buf := bytes.Buffer{}
		if err := t.Execute(&buf, struct{
			CSS template.CSS
			JS template.JS
		}{indexCSS, indexJS}); err != nil {
			panic(err)
		}
		indexHTMLRendered = buf.Bytes()
	}
	if env.Config.Env == "production" {
		once.Do(work)
		return
	}
	work()
}

func getChartTemplate() *template.Template {
	compile()
	return chartHTMLCompiled
}

func getIndexRendered() []byte {
	compile()
	return indexHTMLRendered
}

func renderChart(w http.ResponseWriter, chart *Chart) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	if err := getChartTemplate().Execute(w, struct{
		Chart *Chart
		CSS template.CSS
		JS template.JS
	}{chart, indexCSS, indexJS}); err != nil {
		panic(err)
	}
}