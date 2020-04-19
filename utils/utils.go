package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

// GetProtocol ...
func GetProtocol(r *http.Request) string {
	if r.Header.Get("X-Forwarded-Proto") == "https" || r.TLS != nil {
		return "https"
	}
	return "http"
}

// SetBasicAuth ...
func SetBasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// HTTPClientHeader Headers for HTTP Client
type HTTPClientHeader struct {
	Key   string
	Value string
}

// HTTPClientResponse ...
type HTTPClientResponse struct {
	Body       []byte
	Cookie     []*http.Cookie
	Header     map[string][]string
	StatusCode int
}

// HTTPClientReq ...
func HTTPClientReq(
	clientURL string,
	postParams []byte,
	reqHeaders []*HTTPClientHeader,
	reqCookies []*http.Cookie,
) (*HTTPClientResponse, error) {

	req, err := http.NewRequest("GET", clientURL, nil)

	if postParams != nil {
		var postParamsJSON = []byte(string(postParams))
		req, err = http.NewRequest("POST", clientURL, bytes.NewBuffer(postParamsJSON))
	}

	// Pass received headers
	if len(reqHeaders) > 0 {
		for _, v := range reqHeaders {
			req.Header.Set(v.Key, v.Value)
		}
	}

	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return &HTTPClientResponse{}, err
	}

	// Pass received cookies
	cookieJar.SetCookies(req.URL, reqCookies)

	client := &http.Client{
		Jar: cookieJar,
	}
	resp, err := client.Do(req)
	if err != nil {
		return &HTTPClientResponse{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// Set Cookies
	var httpCookie []*http.Cookie

	secure := false
	if GetProtocol(req) == "https" {
		secure = true
	}

	for _, v := range cookieJar.Cookies(req.URL) {
		httpCookie = append(httpCookie, &http.Cookie{
			Name:     v.Name,
			Value:    v.Value,
			Path:     v.Path,
			Domain:   v.Domain,
			Secure:   secure,
			HttpOnly: v.HttpOnly,
		})
	}

	// Set Headers
	httpHeaders := make(map[string][]string)

	if resp.Header != nil {
		for k, v := range resp.Header {
			switch k {
			case "Set-Cookie":

			default:
				httpHeaders[k] = v
			}

		}
	}

	return &HTTPClientResponse{
		Body:       body,
		Cookie:     httpCookie,
		Header:     httpHeaders,
		StatusCode: resp.StatusCode,
	}, nil
}

// GetBytes receives and interface and returns an array of bytes
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
