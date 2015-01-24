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

type HtmlValueContainer struct {
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
	HtmlValues []HtmlValueContainer
}

type Config struct {
	Init    bool
	DocPath string
}

func main() {
	firstApi := HtmlValueContainer{MethodType: "GET", CurrentPath: "/login/:id", RequestHeader: map[string]string{"Content-Type": "application/json", "Accept": "application/json"},
		RequestBody: "{ 'main' : { 'id' : 2, 'name' : 'Gopher' }}"}

	secondApi := HtmlValueContainer{MethodType: "POST", CurrentPath: "/singup", RequestHeader: map[string]string{"Content-Type": "application/json", "Accept": "application/json"},
		ResponseBody: "{ 'main' : { 'Key' : 'ABC-123-XYZ', 'name' : 'Gopher' }}"}

	config := Config{Init: false, DocPath: "/home/Kaustubh/Desktop/go-projects/yaag/src/yaag/yaag/html/home.html"}

	valueArray := []HtmlValueContainer{secondApi, firstApi}
	allApis := ApiCallValue{BaseLink: "www.google.com", HtmlValues: valueArray}
	GenerateHtml(&allApis, &config)
}

func GenerateHtml(value *ApiCallValue, config *Config) {
	t := template.New("API Documentation")
	filePath, err := filepath.Abs(config.DocPath)
	file, err := ioutil.ReadFile("homeTemplate.html")
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

	/*for _, apiCallValue := range value.HtmlValues {
		t.Execute(homeWriter, map[string]interface{}{"MethodType": apiCallValue.MethodType, "BaseLink": value.BaseLink, "CurrentPath": apiCallValue.CurrentPath, "RequestHeaders": apiCallValue.RequestHeader,
			"RequestUrlParams": apiCallValue.RequestUrlParams, "RequestBody": apiCallValue.RequestBody,
			"ResponseBody": apiCallValue.ResponseBody})
	}*/

	t.Execute(homeWriter, map[string]interface{}{"array": value.HtmlValues,
		"BaseLink": value.BaseLink})

}
