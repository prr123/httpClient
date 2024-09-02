// httpClientV2
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
	"bufio"
	"strings"
	"time"
    "net/http"


 	util "github.com/prr123/utility/utilLib"
)

func main() {

    numarg := len(os.Args)
    dbg := false
//	addrStr :=""
//	msgStr := ""
//	methStr := "GET"
//	cliStr := ""

    flags:=[]string{"dbg","dial"}

//    useStr := "httpClient /dial=adr [/method=methStr] [/cli=clicmd] [/msg=msgStr] [/timing][/dbg]"
    useStr := "httpClient /dial=adr [/dbg]"
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

	timSw:=true
    _, ok = flagMap["timing"]
    if !ok {timSw = }
*/
	addrStr:=""
    val, ok := flagMap["dial"]
    if !ok {
        log.Fatalf(" error no dest provided!\n")
    } else {
        if val.(string) == "none" {log.Fatalf("error: no adr string provided!\n")}
        addrStr = val.(string)
    }

	if dbg {
		fmt.Println("****** setup values ******")
    	fmt.Printf("  debug:  %t\n", dbg)
		fmt.Printf("  dial:   %s\n", addrStr)
		fmt.Println("******* end setup ********")
	}

	// todo verify address string
//    url := "http://89.116.30.49:12001"
	dest := "http:/" + "/" + addrStr +"/"


    // Create a Bearer string by appending string access token
//    bearer := "Bearer " + tokStr


    // add authorization header to the req
//    req.Header.Add("Authorization", bearer)

    // Send req using http Client
    client := &http.Client{}

	// synchronous operation: client is blocked until response is received
	var timSt time.Time
	var timElaps, timElaps1 time.Duration
	var resp *http.Response

    reader := bufio.NewReader(os.Stdin)

    for {
		fmt.Printf("method: x[ help],g,p,q,d,h,c,o,t,a>>")
        inpMeth, _ := reader.ReadString('\n')
		methStr :=""
		switch inpMeth[0] {
			case 'g':
				methStr = "GET"
			case 'p':
				methStr = "POST"
			case 'q':
				methStr = "PUT"
			case 'd':
				methStr = "DELETE"
			case 'h':
				methStr = "HEAD"
			case 'c':
				methStr = "CONNECT"
			case 'o':
				methStr = "OPTIONS"
			case 't':
				methStr = "TRACE"
			case 'a':
				methStr = "PATCH"
			case  'x':
				fmt.Printf("help: valid inputs are >> x [help], g [get], p [post], q [put], d [delete], h [head], c [connect], o [options], t [trace], a [patch]\n")
				continue
			default:
				fmt.Printf("invalid method: %q\n", inpMeth[0])
				fmt.Printf("valid inputs are>> x [help], g [get], p [post], q [put], d [delete], h [head], c [connect], o [options], t [trace], a [patch]\n")
				continue
		}
		fmt.Printf("http cli: >>")
        cliStr, _ := reader.ReadString('\n')
        fmt.Print("body>> ")
        bodyStr, _ := reader.ReadString('\n')
//        fmt.Fprintf(con, bodyStr + "\n")
 	   // Create a new request using http
		if dbg {
			fmt.Printf("method: %s\n", methStr)
			fmt.Printf("cli: %s", cliStr)
			fmt.Printf("body:   %s", bodyStr)
		}

		bodyReader := bytes.NewReader([]byte(bodyStr[:len(bodyStr)-1]))

		cliStr = string(cliStr[:len(cliStr) -1])
		url:= dest + cliStr
		if dbg {fmt.Printf("method: %s url: %s\n", methStr, url)}
		switch methStr {
		case "GET":
			timSt = time.Now()
			resp, err = client.Get(url)
			if err != nil {
				log.Printf("Error GET response: %v\n", err)
				continue
			}
		case "POST":
			timSt = time.Now()
			resp, err = client.Post(url, "text/plain", bodyReader)
			if err != nil {
				log.Printf("Error Post response: %v\n", err)
				continue
			}

		case "HEAD":
			timSt = time.Now()
			resp, err = client.Head(url)
			if err != nil {
				log.Printf("Error Head response: %v\n", err)
				continue
			}


		default:
			req, err := http.NewRequest(methStr, url, bodyReader)
			if err != nil {
				log.Printf("New Req: %v\n", err)
				continue
			}

			timSt = time.Now()
		    resp, err = client.Do(req)
        	if strings.TrimSpace(bodyStr) == "STOP" {
            	log.Println("http client exiting...")
            	os.Exit(0)
        	}

    		if err != nil {
        		log.Printf("Error on response: %v\n", err)
				continue
    		}
		}

		log.Printf("resp status: %s code: %d\n", resp.Status, resp.StatusCode)
		timElaps1 = time.Since(timSt)
//		defer resp.Body.Close()
	    body, err := ioutil.ReadAll(resp.Body)
    	if err != nil {
        	log.Printf("Error while reading the response bytes: %v", err)
    	}
		resp.Body.Close()

		timElaps = time.Since(timSt)
    	log.Printf("resp [rt: %4.1fms %4.1fms %d]: %s\n", float64(timElaps1/time.Millisecond), float64(timElaps/time.Millisecond), len(body), string(body))

	}
}
