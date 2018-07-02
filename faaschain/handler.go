package function

import (
	"fmt"
	"handler/function/sdk"
	"handler/function/store"
	"net/http"
	"os"
	"strings"
)

func makeQueryStringFromParam(params map[string]string) string {
	if params == nil {
		return ""
	}
	result := ""
	for key, value := range params {
		keyVal := fmt.Sprintf("%s-%s", key, value)
		if result == "" {
			result = "?" + keyVal
		} else {
			result = result + "&" + keyVal
		}
	}
	return result
}

func buildUpstreamRequest(function string, data []byte, param map[string]string) *http.Request {
	url := "http://" + function + ":8080"
	queryString := makeQueryStringFromParam(param)
	if queryString != "" {
		url = url + queryString
	}

	var method string

	if method, ok := param["method"]; !ok {
		method = os.Getenv("default-method")
		if method == "" {
			method = "POST"
		}
	}

	req, _ := http.NewRequest(method, deviceUrl, bytes.NewBuffer(data))

}

func execute(request *Request) string {
	var def *Request

	// if store is not defined
	if os.Getenv("store-url") == "" {
		def = request
	} else {
		def, err := store.GetChain(request.Name)
		if err != nil {
			log.Printf("failed to get chain from store, error %v", err)
			return fmt.Sprintf("failed to get chain from store, error %v", err)
		}
	}

	var result string

	// Execute all function
	for index, execute := range def.Executes {
		function := execute.Name
		params := execute.Params
		req := buildUpstreamRequest(function, request.Body.Raw, params)
		client := &http.Client{}

	}

}

// Handle a serverless request
func Handle(req []byte) string {
	request, err := sdk.ParseRequest(req)
	if err != nil {
		log.Printf("failed to parse request object, error %v", err)
		return fmt.Printf("failed to parse request object, error %v", err)
	}

	switch request.Type {
	case sdk.EXECUTE:
		return execute(request)

	case sdk.DEFINE:
		return define(request)

	case sdk.REMOVE:
		return remove(request)

	default:
		log.Printf("invalid request type received '%s'", request.Type)
		return fmt.Printf("failed to parse request object, error %v", err)
	}
}
