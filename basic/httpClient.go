// httpClient
//
// author: prr, azul software
// date: 3/11/2023
// copyright (c) 2023 prr, azulsoftware
//
// simple httpClient for testing static servers
//

package main

import (
	"os"
	"fmt"
    "io/ioutil"
    "log"
	"bytes"
	"time"
    "net/http"


 	util "github.com/prr123/utility/utilLib"
)

func main() {

    numarg := len(os.Args)
    dbg := false
	addrStr :=""
	msgStr := ""
	methStr := "GET"
	cliStr := ""

    flags:=[]string{"dbg","dial", "method", "cli", "msg", "timing"}

    useStr := "httpClient /dial=adr [/method=methStr] [/cli=clicmd] [/msg=msgStr] [/timing][/dbg]"
    helpStr := "simple http Client\n"

    if numarg > len(flags) +1 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: %s", useStr)
        os.Exit(-1)
    }

    if numarg > 1 && os.Args[1] == "help" {
        fmt.Printf("help: %s\n", helpStr)
        fmt.Printf("usage is: %s\n", useStr)
        os.Exit(1)
    }

    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}

    _, ok := flagMap["dbg"]
    if ok {dbg = true}
/*
    if dbg {
        fmt.Printf("dbg -- flag list:\n")
        for k, v :=range flagMap {
            fmt.Printf("  flag: /%s value: %s\n", k, v)
        }
    }
*/
	timSw:=false
    _, ok = flagMap["timing"]
    if ok {timSw = true}

    val, ok := flagMap["dial"]
    if !ok {
        log.Fatalf(" error no dest provided!\n")
    } else {
        if val.(string) == "none" {log.Fatalf("error: no adr string provided!\n")}
        addrStr = val.(string)
    }

    methval, ok := flagMap["method"]
    if ok {
        if val.(string) == "none" {log.Fatalf("error: no method string provided!\n")}
        methStr = methval.(string)
    }

    clival, ok := flagMap["cli"]
    if ok {
        if clival.(string) == "none" {log.Fatalf("error: no cli string provided!\n")}
        cliStr = clival.(string)
    }

    msgval, ok := flagMap["msg"]
    if ok {
        if msgval.(string) == "none" {log.Fatalf("error: no msg string provided!\n")}
        msgStr = msgval.(string)
    }

	if dbg {
		fmt.Println("****** setup values ******")
    	fmt.Printf("  debug:  %t\n", dbg)
		fmt.Printf("  dial:   %s\n", addrStr)
		fmt.Printf("  method: %s\n", methStr)
		fmt.Printf("  cli:    %s\n", cliStr)
		fmt.Printf("  msg:    %s\n", msgStr)
		fmt.Printf("  timing: %t\n", timSw)
		fmt.Println("******* end setup ********")
	}

	// todo verify address string
//    url := "http://89.116.30.49:12001"
	url := `http://` + addrStr +"/" + cliStr


    // Create a Bearer string by appending string access token
//    bearer := "Bearer " + tokStr

    // Create a new request using http
	bodyReader := bytes.NewReader([]byte(msgStr))

	req, err := http.NewRequest(methStr, url, bodyReader)

    // add authorization header to the req
//    req.Header.Add("Authorization", bearer)

    // Send req using http Client
    client := &http.Client{}

	// synchronous operation: client is blocked until response is received
	var timSt time.Time
	var timElaps time.Duration
	if timSw {timSt = time.Now()}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Error on response: %v\n", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error while reading the response bytes: %v", err)
    }
	if timSw {timElaps = time.Since(timSt)}
    log.Printf("resp [ms: %5.1f] body [%d]: %s\n", float64(timElaps/time.Millisecond), len(body), string(body))
}
