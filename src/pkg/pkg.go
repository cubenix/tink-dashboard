package pkg

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	errDirectory     = "failed to open directory"
	errReadContent   = "failed to read directory content"
	errTemplate      = "failed to read template"
	errTemplateParse = "failed to parse template"
)

// PopulateTemplates reads and parses all the available templates
func PopulateTemplates() map[string]*template.Template {
	pwd, _ := os.Getwd()
	basePath := join(pwd, "src", "app", "templates")

	layout := template.Must(
		template.ParseFiles(join(basePath, "_layout.html")),
	)
	template.Must(
		layout.ParseFiles(
			join(basePath, "_header.html"),
			join(basePath, "_footer.html"),
		),
	)
	dir, err := os.Open(join(basePath, "content"))
	CheckError(err, errDirectory)

	fis, err := dir.Readdir(-1)
	CheckError(err, errReadContent)

	result := map[string]*template.Template{}
	for _, fi := range fis {
		f, err := os.Open(join(basePath, "content", fi.Name()))
		CheckError(err, errTemplate+" "+fi.Name())
		content, err := ioutil.ReadAll(f)
		CheckError(err, errTemplate+" "+fi.Name())
		f.Close()

		tmpl := template.Must(layout.Clone())
		_, err = tmpl.Parse(string(content))
		CheckError(err, errTemplateParse+" "+fi.Name())
		result[fi.Name()] = tmpl
	}
	return result
}

// CheckError checks if there is an error.
// If so, prefix the error with given message.
func CheckError(err error, message string) {
	if err != nil {
		log.Fatalf("%v: %v", message, err)
	}
}

func join(paths ...string) string {
	return filepath.Join(paths...)
}
