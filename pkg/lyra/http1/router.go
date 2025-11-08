package http1

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

func (r *Router) GetResponseFunc(req *Request) (bool, HandleFunc) {
	if val, ok := r.router[req.GetMethod()]; ok {
		if h, ok := val[req.GetPath()]; ok {
			return true, h
		}
		return false, NotFound
	}
	return false, MethodNotAllowed
}

func (r *Router) GetResponse(request *Request) *Response {
	_, resFunc := r.GetResponseFunc(request)
	response := resFunc(request)
	return response
}
