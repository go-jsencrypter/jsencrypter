package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"log"
	"net/http"
	"net/url"
)

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() && strings.Contains(f.Name(), ".js") {
		log.Printf("Encrypt: %s\n", path)
		body, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("An error occurred while read file:%s,error:%v\n", f.Name(), err)
		}

		resp, err := http.PostForm("http://tool.lu/js/ajax.html", url.Values{"code": {string(body)}, "operate": {"pack"}})
		if err != nil {
			log.Printf("An error occurred while connecting server to encrypt file:%s, error: %v\n", f.Name(), err)
		}
		defer resp.Body.Close()
		rBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("An error occurred while read response body to encrypt file: %s, error: %v \n", f.Name(), err)
		}

		//encrypted js file to out dir
		dest := filepath.Join("out", path)

		resData := make(map[string]string)
		json.Unmarshal(rBody, &resData)
		var pathReg string
		if string(os.PathSeparator) == "\\" {
			pathReg = "\\\\([\\w]*).js"
		} else {
			pathReg = "/([\\w]*).js"
		}
		re := regexp.MustCompile(pathReg)
		outDir := re.ReplaceAllLiteralString(dest, "")
		err = os.MkdirAll(outDir, 0777)
		if err != nil {
			log.Printf("An error occurred while mkdir,file: %s,error: %v\n", f.Name(), err)

		}

		//write file from http response data to local file(dest)
		if err := ioutil.WriteFile(dest, []byte(resData["text"]), 0666); err != nil {
			log.Printf("An error occurred while write file to dest,file: %s,error: %v\n", f.Name(), err)
		}
	}

	return nil
}

func main() {
	flag.Parse()
	root := flag.Arg(0)
	err := filepath.Walk(root, visit)
	if (err != nil) {
		log.Println("encrypt js file failed.")
		os.Exit(-1);
	}
	log.Println("Encrypt all js file successed.")
}
