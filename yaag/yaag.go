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
	_ "path/filepath"
)

type APICall struct {
	BaseLink    string
	CurrentPath string
	MethodType  string

	PostForm map[string]string

	RequestHeader    map[string]string
	ResponseHeader   map[string]string
	RequestUrlParams map[string]string

	RequestBody  string
	ResponseBody string
}

func main() {
	value := APICall{BaseLink: " http://www.facebook.com ",
		MethodType: "GET", CurrentPath: "/login/:id", RequestHeader: map[string]string{"Content-Type": "application/json", "Accept": "application/json"},
		RequestBody: "{ 'main' : { 'id' : 2, 'name' : 'Gopher' }}"}
	GenerateHtml(&value)
}

func GenerateHtml(value *APICall) {
	t := template.New("API Documentation")
	/*filePath, err := filepath.Abs("../htmlTemplate.html")
	log.Println(filePath)*/
	file, err := ioutil.ReadFile("homeTemplate.html")
	if err != nil {
		log.Println(err)
		return
	}
	htmlString := string(file)
	//log.Println("Html String ", htmlString)
	t, err = t.Parse(htmlString)
	if err != nil {
		log.Println(err)
		return
	}
	homeHtmlFile, err := os.Create("home.html")
	defer homeHtmlFile.Close()

	if err != nil {
		log.Println(err)
		return
	}
	homeWriter := io.Writer(homeHtmlFile)

	t.Execute(homeWriter, map[string]interface{}{"MethodType": value.MethodType, "BaseLink": value.BaseLink, "CurrentPath": value.CurrentPath, "RequestHeaders": value.RequestHeader,
		"RequestUrlParams": value.RequestUrlParams, "RequestBody": value.RequestBody})

}
