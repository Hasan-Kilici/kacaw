package kacaw

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
  "os"
  "io"
  "mime/multipart"
)

type Router struct {
	Routes map[string]map[string]http.HandlerFunc
}

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func Default() *Router {
	return &Router{
		Routes: make(map[string]map[string]http.HandlerFunc),
	}
}

func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.registerHandler("GET", path, handler)
}

func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.registerHandler("POST", path, handler)
}

func (r *Router) Run(port string) {
	http.ListenAndServe(port, r)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	if route, ok := r.Routes[method][path]; ok {
		context := &Context{
			ResponseWriter: w,
			Request:        req,
		}
		route(context.ResponseWriter, context.Request)
	} else {
		http.NotFound(w, req)
	}
}

func HTML(status int, filename string, data interface{}, w http.ResponseWriter) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
}

func JSON(status int, data interface{}, w http.ResponseWriter) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
}

func (r *Router) registerHandler(method, path string, handler http.HandlerFunc) {
	if _, ok := r.Routes[method]; !ok {
		r.Routes[method] = make(map[string]http.HandlerFunc)
	}
	r.Routes[method][path] = handler
}

func (r *Router) LoadHTMLFiles(filenames ...string) {
	templates := []string{}
	for _, filename := range filenames {
		matches, err := filepath.Glob(filename)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		templates = append(templates, matches...)
	}

	for _, templateFile := range templates {
		tmpl, err := template.ParseFiles(templateFile)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		filename := filepath.Base(templateFile)
		r.GET("/"+filename, func(w http.ResponseWriter, req *http.Request) {
			tmpl.Execute(w, nil)
		})
	}
}

func (r *Router) Static(path, directory string) {
	fs := http.StripPrefix(path, http.FileServer(http.Dir(directory)))
	r.GET(path+"/{filepath:*}", fs.ServeHTTP)
}

func Redirect(w http.ResponseWriter, req *http.Request, url string) {
	http.Redirect(w, req, url, http.StatusFound)
}

func SaveFile(file multipart.File, header *multipart.FileHeader, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}

	return nil
}
