/*
 * This is the main core of the yaag package
 */
package yaag

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type APICall struct {
	Id int

	CurrentPath string
	MethodType  string

	PostForm map[string]string

	RequestHeader    map[string]string
	ResponseHeader   map[string]string
	RequestUrlParams map[string]string

	RequestBody  string
	ResponseBody string
	ResponseCode int
}

type ApiCallValue struct {
	BaseLink   string
	HtmlValues []APICall
}

type Config struct {
	Init    bool
	DocPath string
}

func main() {

	firstApi := APICall{Id: 1, MethodType: "GET", CurrentPath: "/login/:id", RequestHeader: map[string]string{"Content-Type": "application/json", "Accept": "application/json"},

		RequestBody: "{ 'main' : { 'id' : 2, 'name' : 'Gopher' }}"}

	secondApi := APICall{Id: 2, MethodType: "POST", CurrentPath: "/singup", RequestHeader: map[string]string{"Content-Type": "application/json", "Accept": "application/json"},
		ResponseBody: "{ 'main' : { 'Key' : 'ABC-123-XYZ', 'name' : 'Gopher' }}"}

	config := Config{Init: false, DocPath: "html/home.html"}

	valueArray := []APICall{secondApi, firstApi}
	allApis := ApiCallValue{BaseLink: "www.google.com", HtmlValues: valueArray}
	GenerateHtml(&allApis, &config)
}

func GenerateHtml(value *ApiCallValue, config *Config) {
	t := template.New("API Documentation")
	filePath, err := filepath.Abs(config.DocPath)
	file, err := ioutil.ReadFile("templates/main.html")
	if err != nil {
		log.Println(err)
		return
	}
	htmlString := string(file)
	t, err = t.Parse(htmlString)
	if err != nil {
		log.Println(err)
		return
	}
	homeHtmlFile, err := os.Create(filePath)
	defer homeHtmlFile.Close()

	if err != nil {
		log.Println(err)
		return
	}
	homeWriter := io.Writer(homeHtmlFile)
	t.Execute(homeWriter, map[string]interface{}{"array": value.HtmlValues,
		"BaseLink": value.BaseLink})

}
