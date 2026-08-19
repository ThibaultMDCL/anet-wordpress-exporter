package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// TODO: rename ces nom d'objet instacié il faisait quoi ce dev ????

func (c *WordpressCollector) FetchJSONFromEndpoint(APIEndpoint string) ([]byte, error) {
	APIBase := c.Wp.MonitoredWordpress
	HTTPClient := &http.Client{}
	fetchURL := fmt.Sprintf("%s%s", APIBase, APIEndpoint)
	fmt.Println(fetchURL)

	request, err := http.NewRequest(http.MethodGet, fetchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	request.Header.Set("User-Agent", c.Wp.UserAgent)

	if c.Wp.Auth.Use {
		request.SetBasicAuth(
			c.Wp.Auth.Username,
			c.Wp.Auth.Password,
		)
	}

	response, err := HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request %s: %w", fetchURL, err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf(
			"request %s returned HTTP status %d",
			fetchURL,
			response.StatusCode,
		)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return data, nil
}

// count items returned in JSON and return length
func CountJSONItems(JSONResponse []byte) (int, error) {
	var err error
	var JSONObject interface{}
	json.Unmarshal(JSONResponse, &JSONObject)

	JSONObjectSlice, isOK := JSONObject.([]interface{})
	if !isOK {
		err = fmt.Errorf("Cannot convert the JSON object")
		// return -1 if json cannot be parsed properly
		return -1, err
	}

	return len(JSONObjectSlice), err
}

func BasicAuth(username, password string) string {
	authString := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(authString))
}

func ErrCheck(e error) {
	if e != nil {
		log.Println(e)
	}
}
