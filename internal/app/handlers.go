package app

import "net/http"

type MiscController struct{}

func (self *MiscController) test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
