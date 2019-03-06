package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	// send request to /register with token. use http.Request
	url := "https://localhost:8080/register"

	//log.Printf(cookies[0].String())
	uuid, otp := callAPI(url, nil)
	url = "https://localhost:8080/verify_phone"
	var cookies = []http.Cookie{
		http.Cookie{
			Name:    "UUID",
			Value:   uuid,
			Expires: time.Now().AddDate(0, 0, 1),
		},
		http.Cookie{
			Name:    "OTP",
			Value:   otp,
			Expires: time.Now().AddDate(0, 0, 1),
		},
		http.Cookie{
			Name:    "phone_number",
			Value:   "57352919",
			Expires: time.Now().AddDate(0, 0, 1),
		},
	}

	callAPI(url, cookies)
}

func callAPI(url string, cookies []http.Cookie) (string, string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	payload := map[string]string{"first_name": "kaleem", "middle_name": "mohammad", "last_name": "peeroo", "email": "kpeeroo@gmail.com", "phone_number": "57352919"}
	payloadByte, err := json.Marshal(payload)
	fmt.Println(err)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadByte))
	req.Header.Add("content-type", "application/json")
	value := "Bearer " + "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6Ik5UTkJOa1UwT1RoRFFrWTVORFZGUlRBNU1rRkRORUkwTmtNNE1FSkZRa1k1UmpWRE5UVTVRUSJ9.eyJpc3MiOiJodHRwczovL2Rldi13aGVyZXNteWRvc2guYXV0aDAuY29tLyIsInN1YiI6IjlmTHZkMzdiMnM2Rmk4WHpCRjFoRkZDY1hsb3JNdGtlQGNsaWVudHMiLCJhdWQiOiJ3aGVyZXNteWRvc2giLCJpYXQiOjE1NTE0NzM0NzYsImV4cCI6MTU1MTU1OTg3NiwiYXpwIjoiOWZMdmQzN2IyczZGaThYekJGMWhGRkNjWGxvck10a2UiLCJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMifQ.gpfdktwWa9ErGT7FO1EdBnS9iDLLBKG37k6faCugfxmlLKUkFX_nD_boubOPOLW6-WCzY9TctWt67s-sj4IcQ1d5AhFYwF9Q2aLJdFYlBw0aDn4Q2o3k5BRQgMcjj9TEWYfihhDx7MEgMgtymH4DJXU_-YA9vqCkyHHrHVLVAlLb9WiLoKZtdtJGW4zUjiGke2PGsh64jwrWTagW8R4iHySwbmmAo_PEgi0lgSge99aE6DgTUAytVSZ9j3EhTtsFDk8u5CGkMRvOow5_OEkbze-RQTp_R58mxIxIWiIOVFVWm-ZtnjfQ8qogbzdTWCxpLzzHO9QFVcZLLrIRCAtgRw"
	req.Header.Add("authorization", value)

	if cookies != nil {
		req.Header.Add("cookie", cookies[0].String())
		req.Header.Add("cookie", cookies[1].String())
		req.Header.Add("cookie", cookies[2].String())
		fmt.Println(req)
	}
	// send req get response from server
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	cookies_ := res.Cookies()

	// test cookie values
	var uuid, otp string
	for _, cookie := range cookies_ {
		//jwtsMap[jwt.UUIDKey.String()] = requestUUID
		//jwtsMap[jwt.OTP] = otp
		fmt.Println("looking for cookies ...")
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
		if cookie.Name == "UUID" {
			uuid = cookie.Value
		}
		if cookie.Name == "OTP" {
			otp = cookie.Value
		}
		log.Print(res)
		//fmt.Println(string(body))
		log.Print(body)
	}
	return uuid, otp
}
