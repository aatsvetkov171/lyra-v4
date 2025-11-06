package http1

func NotFound(request *Request) *Response {
	response := NewResponse()
	response.AddBody([]byte("<p>Lyyyyraaa</p>"))
	return response //потом исправлю на правильное поведение
}

func MethodNotAllowed(request *Request) *Response {
	response := NewResponse()
	response.AddBody([]byte("<p>Lyyyyraaa</p>"))
	return response //потом исправлю на правильное поведение
}

func Hello(request *Request) *Response {
	response := NewResponse()
	response.AddBody([]byte(`
	<p>Lyyyyraaa Helloooo</p>
	<p>Lyyyyraaa Helloooo</p>
	<p>Lyyyyraaa Helloooo</p>
	<p>Lyyyyraaa Helloooo</p>
	<p>Lyyyyraaa Helloooo</p>
	<p>Lyyyyraaa Helloooo</p>
	`))
	return response
}

func About(request *Request) *Response {
	response := NewResponse()
	response.AddBody(
		[]byte(`
		<h1>О нас</h1>
		<p>safl kjs kj dkf ksfk;fdk; snfk fdkld lkdf </p>
		<p>afsjjka kjafsk fsjak afk</p>
		<h2>jghf</h2>
		<p>sgdlsdlgjk sgg sd kj gsgskhg hkls ksjk sdgsg d</p>
		`))
	return response
}
