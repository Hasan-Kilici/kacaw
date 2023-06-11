package main

import (
	"encoding/json"
  	"html/template"
	"net/http"
	"time"
)

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	staticPath     string
}

type Router struct {
	Routes        map[string]map[string]http.HandlerFunc
	staticPath    string
	CookieManager CookieManager
}

type RouterOptions struct {
	CookieManager CookieManager
}

func Default(options ...RouterOptions) *Router {
	r := &Router{
		Routes: make(map[string]map[string]http.HandlerFunc),
	}

	if len(options) > 0 {
		cookieManager := options[0].CookieManager
		if cookieManager == nil {
			cookieManager = &DefaultCookieManager{}
		}
		r.CookieManager = cookieManager
	} else {
		r.CookieManager = &DefaultCookieManager{}
	}

	return r
}

func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.registerHandler("GET", path, handler)
}

func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.registerHandler("POST", path, handler)
}

func (r *Router) HEAD(path string, handler http.HandlerFunc) {
	r.registerHandler("HEAD", path, handler)
}

func (r *Router) PUT(path string, handler http.HandlerFunc) {
	r.registerHandler("PUT", path, handler)
}

func (r *Router) DELETE(path string, handler http.HandlerFunc) {
	r.registerHandler("DELETE", path, handler)
}

func (r *Router) CONNECT(path string, handler http.HandlerFunc) {
	r.registerHandler("CONNECT", path, handler)
}

func (r *Router) OPTIONS(path string, handler http.HandlerFunc) {
	r.registerHandler("OPTIONS", path, handler)
}

func (r *Router) TRACE(path string, handler http.HandlerFunc) {
	r.registerHandler("TRACE", path, handler)
}

func (r *Router) PATCH(path string, handler http.HandlerFunc) {
	r.registerHandler("PATCH", path, handler)
}

func (r *Router) Run(port string) {
	server := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	server.SetKeepAlivesEnabled(true)

	r.enableCaching()

	server.ListenAndServe()
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	if route, ok := r.Routes[method][path]; ok {
		context := &Context{
			ResponseWriter: w,
			Request:        req,
			staticPath:     r.staticPath,
		}
		route(context.ResponseWriter, context.Request)
	} else {
		http.NotFound(w, req)
	}
}

func (r *Router) HTML(status int, filename string, data interface{}, w http.ResponseWriter) {
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

func (r *Router) JSON(status int, data interface{}, w http.ResponseWriter) {
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
	if r.Routes[method] == nil {
		r.Routes[method] = make(map[string]http.HandlerFunc)
	}

	r.Routes[method][path] = handler
}

func (r *Router) SaveFile(file multipart.File, header *multipart.FileHeader, path string) error {
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

func (r *Router) enableCaching() {
	if r.staticPath != "" {
		fs := http.FileServer(http.Dir(r.staticPath))
		r.GET("/static/*filepath", func(w http.ResponseWriter, req *http.Request) {
			req.URL.Path = req.URL.Path[len("/static"):]
			fs.ServeHTTP(w, req)
		})
	}
}
