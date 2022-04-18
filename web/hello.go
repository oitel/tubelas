package web

import "net/http"

const (
	helloHtml = `
<html>
	<h1>Hello, world!</h1>
</html>
`
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(helloHtml))
}
