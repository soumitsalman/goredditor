package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

// type to contain http errors when the status code is NOT 200
type HttpStatusCodeError int

func (err *HttpStatusCodeError) Error() string {
	return fmt.Sprintf("Request failed with status code: %v", *err)
}

func IsAuthTokenExpired(auth_token string) bool {
	var jwtToken, _, err = new(jwt.Parser).ParseUnverified(auth_token, jwt.MapClaims{})
	if err != nil {
		return true
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return true
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return true
	}
	return float64(time.Now().Unix()) > exp
}

func MakeBearerToken(bearer_token string) string {
	return "bearer " + bearer_token
}

func MakeBasicAuthToken(username string, password string) string {

	//authValue = fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(client.creds.ClientId+":"+client.creds.ClientSecret)))

	return fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
}

// this sends and deserialized the payload assuming the response if json blob
func GetHttpResponse(req *http.Request, http_client *http.Client) (map[string]any, error) {

	if resp, err := http_client.Do(req); err != nil {
		//log.Println("Getting New Access Token Failed: ", err)
		return nil, err
	} else if resp.StatusCode == http.StatusUnauthorized {
		//log.Println("Getting New Access Token Failed: ")
		unauthErr := HttpStatusCodeError(resp.StatusCode)
		return nil, &unauthErr
	} else {
		defer resp.Body.Close()
		return DeserialzeJsonBlob[map[string]any](io.Reader(resp.Body))
	}
}

// this is a shell http request builder and does NOT actually verify if the access token is valid
func BuildHttpRequest(method string, url string, payload io.Reader, payload_encoding string, auth_token string, user_agent string, ctx *context.Context) *http.Request {
	req, _ := http.NewRequestWithContext(*ctx, method, url, payload)
	//standard header
	req.Header.Add("Authorization", auth_token)
	req.Header.Add("User-Agent", user_agent)
	//assign content type based on encoding
	if payload != nil {
		switch payload_encoding {
		case "json":
			req.Header.Add("Content-Type", "application/json; charset=UTF-8")
		case "url":
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		default:
			log.Printf("unsupport/unrecognized encoding: %v\n", payload_encoding)
			return nil
		}
	}
	return req
}

/*
// this is a shell http request builder and does NOT actually verify if the access token is valid
func BuildHttpRequest(method string, url string, payload io.Reader, payload_encoding string, auth_token string, user_agent string, ctx *context.Context) *http.Request {
	var req, _ = http.NewRequestWithContext(*ctx, method, url, payload)
	//standard header
	req.Header.Add("Authorization", auth_token)
	req.Header.Add("User-Agent", user_agent)
	//assign content type based on encoding
	if payload != nil {
		switch payload_encoding {
		case "json":
			req.Header.Add("Content-Type", "application/json; charset=UTF-8")
		case "url":
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		default:
			log.Printf("unsupport/unrecognized encoding: %v\n", payload_encoding)
			return nil
	}
	return req
}
*/
