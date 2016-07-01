package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"net/http"
	"net/url"
)

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	fmt.Printf("isDir:%v\n", f.IsDir())
	if !f.IsDir() && strings.Contains(f.Name(), ".js") {
		fmt.Printf("isJS:%s\n", f.Name())
		body, err := ioutil.ReadFile(path)
		if err != nil {
			//todo
		}
		fmt.Printf("body:%s\n", body)

		resp, err := http.PostForm("http://tool.lu/js/ajax.html", url.Values{"code": {string(body)}, "operate": {"pack"}})
		if err != nil {
			fmt.Println("http error")
		}
		defer resp.Body.Close()
		rBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("http return error")
		}
		fmt.Printf("return body:%s\n", rBody)
		tmpfn := filepath.Join("out", path)

		fileContent := make(map[string]string)
		json.Unmarshal(rBody, &fileContent)
		fmt.Println("dir:" + tmpfn)
		re := regexp.MustCompile("/([\\w]*).js")
		fmt.Printf("replace:%s\n", re.ReplaceAllLiteralString(tmpfn, ""))
		tmpfnDir := re.ReplaceAllLiteralString(tmpfn, "")
		err = os.MkdirAll(tmpfnDir, 0777)
		if err != nil {
			fmt.Println(err)
		}
		// if _, err := os.Stat(tmpfnDir); os.IsNotExist(err) {
		// 	// path/to/whatever does not exist
		// 	fmt.Println("dir not exsits")
		// 	err = os.Mkdir(tmpfnDir, 0777)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// }
		if err := ioutil.WriteFile(tmpfn, []byte(fileContent["text"]), 0666); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func main() {
	flag.Parse()
	root := flag.Arg(0)
	err := filepath.Walk(root, visit)
	fmt.Printf("filepath.Walk() returned %v\n", err)
}
