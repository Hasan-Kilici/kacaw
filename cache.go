package kacaw

import "net/http"

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
			context := &Context{
				ResponseWriter: w,
				Request:        req,
				staticPath:     r.staticPath,
			}
			tmpl.Execute(w, context)
		})
	}
}

func (r *Router) Static(filenames ...string) {
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
			context := &Context{
				ResponseWriter: w,
				Request:        req,
				staticPath:     r.staticPath,
			}
			tmpl.Execute(w, context)
		})
	}
}

func (r *Router) Redirect(w http.ResponseWriter, req *http.Request, url string) {
	http.Redirect(w, req, url, http.StatusFound)
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
	r.GET("/*filepath", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		http.FileServer(http.Dir(r.staticPath)).ServeHTTP(w, req)
	})
}
