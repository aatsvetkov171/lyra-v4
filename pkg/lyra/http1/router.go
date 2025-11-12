package http1

import (
	"strings"
)

type HandleFunc func(*Request) *Response

var Mime = map[string]map[string]string{
	"media": {
		"jpg":  "image/jpeg; charset=UTF-8",
		"jpeg": "image/jpeg; charset=UTF-8",
		"png":  "image/png; charset=UTF-8",
		"gif":  "image/gif; charset=UTF-8",
		"webp": "image/webp; charset=UTF-8",
		"svg":  "image/svg+xml; charset=UTF-8",
		"ico":  "image/x-icon; charset=UTF-8",
		"pdf":  "application/pdf",
	},
	"static": {
		"css": "text/css; charset=UTF-8",
		"js":  "text/javascript; charset=UTF-8",
	},
	"template": {
		"html": "text/html; charset=UTF-8",
	},
}

type Router struct {
	router    map[string]map[string]HandleFunc
	staticDir string
	mediaDir  string
}

func NewRouter(staicPath string, mediaPath string) *Router {
	newRouter := Router{
		router:    make(map[string]map[string]HandleFunc),
		staticDir: staicPath,
		mediaDir:  mediaPath,
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
		if len(req.GetPath()) >= len(r.staticDir)+2 {
			staticPrefix := "/" + r.staticDir + "/"
			if strings.HasPrefix(req.GetPath(), staticPrefix) {
				return true, r.SendStaticFile
			}
		}
		if len(req.GetPath()) >= len(r.mediaDir)+2 {
			staticPrefix := "/" + r.mediaDir + "/"
			if strings.HasPrefix(req.GetPath(), staticPrefix) {
				return true, r.SendMediaFile
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
	fileEnd := strings.SplitN(request.GetPath(), ".", 2)
	contentType, ok := Mime["static"][fileEnd[1]]
	if ok {
		response.GetHeaders()["Content-Type"] = contentType
	}
	l := len(r.staticDir) + 2
	//fmt.Println(request.GetPath()[l:])
	response.AddFile(request.GetPath()[l:])
	response.AddMime("static")
	return response
}

func (r *Router) SendMediaFile(request *Request) *Response {
	response := NewResponse(200)
	fileEnd := strings.SplitN(request.GetPath(), ".", 2)
	contentType, ok := Mime["media"][fileEnd[1]]
	if ok {
		response.GetHeaders()["Content-Type"] = contentType
	}
	l := len(r.mediaDir) + 2
	response.AddFile(request.GetPath()[l:])
	response.AddMime("media")
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
