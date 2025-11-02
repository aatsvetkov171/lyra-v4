package http1

import (
	"strings"
)

type HandleFunc func(*Request) *Response

type Router struct {
	router map[string]map[string]HandleFunc
}

func NewRouter() *Router {
	newRouter := Router{
		router: make(map[string]map[string]HandleFunc),
	}
	return &newRouter
}

func (r *Router) Handle(method string, path string, h HandleFunc) {
	if r.router[strings.ToUpper(method)] == nil {
		r.router[strings.ToUpper(method)] = make(map[string]HandleFunc)
	}
	r.router[method][path] = h
}

func (r *Router) GetResponse(req *Request) (bool, HandleFunc) {
	if val, ok := r.router[req.GetMethod()]; ok {
		if h, ok := val[req.GetPath()]; ok {
			return true, h
		}
		return false, nil
	}
	return false, nil
}
