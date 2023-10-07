package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/go-xmlfmt/xmlfmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type RequestBody struct {
	Data          string `json:"data"`
	Base64Decode  bool   `json:"base64Decode"`
	PercentDecode bool   `json:"percentDecode"`
	Beautify      bool   `json:"beautify"`
}

type ResponseBody struct {
	Result     string   `json:"result"`
	Operations []string `json:"operations"`
	Error      error    `json:"error"`
}

func processError(c *gin.Context, httpStatusCode int, err error) {
	log.Fatal(err)
	c.IndentedJSON(httpStatusCode, ResponseBody{Error: err})
}

func processData(c *gin.Context) {
	var operations []string
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	response := ResponseBody{}
	if err != nil {
		processError(c, http.StatusBadRequest, err)
		return
	}
	body := RequestBody{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		processError(c, http.StatusBadRequest, err)
		return
	}
	processedData := body.Data

	if body.PercentDecode {
		// Percent decode
		processedData, err = url.QueryUnescape(processedData)
		if err != nil {
			processError(c, http.StatusBadRequest, err)
			return
		}
		operations = append(operations, "percentDecode")
	}

	fmt.Println(processedData)
	// Base64 decode
	if body.Base64Decode {
		base64Decoded, err := b64.StdEncoding.DecodeString(processedData)
		if err != nil {
			processError(c, http.StatusBadRequest, err)
			return
		}
		processedData = string(base64Decoded)
		operations = append(operations, "base64Decode")
	}
	fmt.Println("lol")
	fmt.Println(processedData)

	if body.Beautify {
		// Is valid XML?
		xmlData := new(interface{})
		err = xml.Unmarshal([]byte(processedData), xmlData)
		fmt.Println("lol3")
		fmt.Println(xmlData)
		if err == nil {
			processedData = xmlfmt.FormatXML(processedData, "", "\t")
			operations = append(operations, "xmlBeautify")
		}
	}

	response.Result = processedData
	response.Operations = operations
	fmt.Println("lol2")
	fmt.Println(response.Result)
	fmt.Println(body.Beautify)

	c.IndentedJSON(http.StatusOK, response)
}

func main() {
	router := gin.Default()
	router.Use(limits.RequestSizeLimiter(5000))
	router.POST("/gobble", processData)
	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}
