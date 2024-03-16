package parser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Check if URL is accessible and get Response
func ParseURL(link string) (string, error) {

	// Validate link
	URL, err := url.Parse(link)

	if err != nil || URL.Scheme == "" || URL.Host == "" {
		fmt.Println(err)
		return "", errors.New(`cannot reach URL`)
	}

	// Make a GET request
	resp, err := http.Get(URL.String())

	if err != nil {
		return err.Error(), nil
	}

	// Read response and return to response
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return fmt.Sprintf("%v\n%v\n%v\n%v\n", resp.Status, resp.StatusCode, resp.Proto, string(body)), nil
}
