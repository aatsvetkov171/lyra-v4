package http1

import (
	"strings"
)

type HandleFunc func(*Request) *Response

type Router struct {
	router    map[string]map[string]HandleFunc
	staticDir string
}

func NewRouter(staicPath string) *Router {
	newRouter := Router{
		router:    make(map[string]map[string]HandleFunc),
		staticDir: staicPath,
	}
	return &newRouter
}

func (r *Router) Handle(method string, path string, h HandleFunc) {
	if r.router[method] == nil {
		r.router[method] = make(map[string]HandleFunc)
	}
	r.router[method][path] = h
}

func (r *Router) GET(path string, h HandleFunc) {
	if r.router["GET"] == nil {
		r.router["GET"] = make(map[string]HandleFunc)
	}
	r.router["GET"][path] = h
}

func (r *Router) POST(path string, h HandleFunc) {
	if r.router["POST"] == nil {
		r.router["POST"] = make(map[string]HandleFunc)
	}
	r.router["POST"][path] = h
}

func (r *Router) GetResponseFunc(req *Request) (bool, HandleFunc) {
	if val, ok := r.router[req.GetMethod()]; ok {
		if h, ok := val[req.GetPath()]; ok {
			return true, h
		}
		if len(req.GetPath()) >= 8 {
			staticPrefix := "/" + r.staticDir + "/"
			if strings.HasPrefix(req.GetPath(), staticPrefix) {
				return true, r.SendStaticFile
			}
		}
		return false, r.NotFound
	}
	return false, r.MethodNotAllowed
}

func (r *Router) GetResponse(request *Request) *Response {
	_, resFunc := r.GetResponseFunc(request)
	response := resFunc(request)
	return response
}

//--------------------------------VIEWS router

func (r *Router) SendStaticFile(request *Request) *Response {
	response := NewResponse(200)
	fileEnd := strings.Split(request.GetPath(), ".")
	if len(fileEnd) == 2 {
		switch fileEnd[1] {
		case "css":
			response.GetHeaders()["Content-Type"] = "text/css; charset=UTF-8"
		case "js":
			response.GetHeaders()["Content-Type"] = "text/javascript; charset=UTF-8"
		}
	}
	l := len(r.staticDir) + 2
	//fmt.Println(request.GetPath()[l:])
	response.AddFile(request.GetPath()[l:])
	return response
}

func (r *Router) NotFound(request *Request) *Response {
	response := NewResponse(404)
	response.AddString("<h1>Not Found</h1>")
	return response
}

func (r *Router) MethodNotAllowed(request *Request) *Response {
	response := NewResponse(405)
	response.AddString("<h1>Method Not Allowed</h1>")
	return response
}
