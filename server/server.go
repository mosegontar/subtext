package server

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mosegontar/underbyte"
)

func home(w http.ResponseWriter, req *http.Request) {
	//url := req.URL
	//queryVals := url.Query()
	//img := queryVals.Get("base64")

	t, err := template.ParseFiles("./server/index.html")
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("Rendering template")
	err = t.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Uh oh!", 500)
	}
}

func parseForm(req *http.Request) (message []byte, imageData []byte, err error) {
	req.ParseMultipartForm(30 << 20)

	rawMessage := req.FormValue("message")
	message = []byte(rawMessage)

	f, _, err := req.FormFile("userImage")
	if err != nil {
		log.Println("no image", err.Error())
		return nil, nil, err
	}
	defer f.Close()
	imageData, err = ioutil.ReadAll(f)

	if err != nil {
		log.Println(err.Error())
		return nil, nil, err
	}

	return
}

func uploadFile(w http.ResponseWriter, req *http.Request) (err error) {
	message, data, err := parseForm(req)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	source := underbyte.SourceImageBytes(data)
	ub := underbyte.NewUnderbyteImage(source)

	ub.EncodeMessage([]byte(message))

	buff := new(bytes.Buffer)
	ub.WriteImage(buff)

	h := sha1.New()
	h.Write(buff.Bytes())
	hs := h.Sum(nil)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%x.png", hs))

	fmt.Fprint(w, buff)
	return
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

	log.Println("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}

}

func main() {
	Start()
}
