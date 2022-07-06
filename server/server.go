package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mosegontar/underbyte"
)

func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "welcome\n")
}

func uploadFile(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(30 << 20)

	f, _, err := req.FormFile("myFile")
	if err != nil {
		log.Println(err.Error())
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)

	source := underbyte.SourceImageBytes(data)
	ub := underbyte.NewUnderbyteImage(source)
	fmt.Fprintf(w, "%v", ub)
}

func upload(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		uploadFile(w, req)
	case "GET":
		fmt.Fprintf(w, "not implemented\n")

	}
}

func Start() {
	http.HandleFunc("/", home)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}

}

func main() {
	Start()
}
