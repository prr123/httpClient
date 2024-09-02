// httpClientV2
//
// author: prr, azul software
// date: 3/11/2023
// copyright (c) 2023 prr, azulsoftware
//
// simple httpClient for testing static servers
//
// V4 add https
//

package main

import (
	"os"
	"fmt"
    "io/ioutil"
    "log"
//	"bytes"
	"bufio"
	"strings"
	"time"
    "net/http"
	"crypto/tls"

 	util "github.com/prr123/utility/utilLib"
)

func main() {

    numarg := len(os.Args)
//	addrStr :=""
//	msgStr := ""
//	methStr := "GET"
//	cliStr := ""

    flags:=[]string{"dbg","dial","https"}

//    useStr := "httpClient /dial=adr [/method=methStr] [/cli=clicmd] [/msg=msgStr] [/timing][/dbg]"
    useStr := "/dial=adr [/https] [/dbg]"
    helpStr := "http client"

    if numarg > len(flags) +1 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: %s %s", os.Args[0], useStr)
        os.Exit(-1)
    }

    if numarg > 1 && os.Args[1] == "help" {
        fmt.Printf("help: %s %s\n", os.Args[0], helpStr)
        fmt.Printf("usage: %s %s\n", os.Args[0], useStr)
        os.Exit(0)
    }

    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}

    dbg := false
    _, ok := flagMap["dbg"]
    if ok {dbg = true}

    httpsFlag := false
	_,ok = flagMap["https"]
	if ok {httpsFlag = true}

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
    	fmt.Printf("  https:  %t\n", httpsFlag)
		fmt.Printf("  dial:   %s\n", addrStr)
		fmt.Println("******* end setup ********")
	}


	dest := "http:/" + "/" + addrStr +"/"

	// need to add https client

/*
	tr := &http.Transport{
    	DisableKeepAlives:   false,
    	MaxIdleConns:        0,
    	MaxIdleConnsPerHost: 0, // no limit
    	IdleConnTimeout:     time.Second * 10,
	}
    client := &http.Client{Transport: tr}
*/

    client := &http.Client{}

	if httpsFlag {
		dest = "https:/" + "/" + addrStr +"/"
		tr := &http.Transport{
        		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    		}
		client = &http.Client{Transport: tr}
	}


	// synchronous operation: client is blocked until response is received
	var timSt time.Time
	var timElaps, timElaps1 time.Duration
	var resp *http.Response

    reader := bufio.NewReader(os.Stdin)

    for {
		fmt.Printf("method: z[help],x,g,p,q,d,h,c,o,t,a>>")
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
			case 'x':
				fmt.Printf("exiting\n")
				os.Exit(0)
			case 'z':
				fmt.Printf("help: valid inputs are >> z[help], x[exit],  g [get], p [post], q [put], d [delete], h [head], c [connect], o [options], t [trace], a [patch]\n")
				continue
			default:
				fmt.Printf("invalid method: %q\n", inpMeth[0])
				fmt.Printf("help: valid inputs are >> z[help], x[exit],  g [get], p [post], q [put], d [delete], h [head], c [connect], o [options], t [trace], a [patch]\n")
				continue
		}
		fmt.Printf("http cli: >>")
        cliStr, _ := reader.ReadString('\n')
        fmt.Print("body>> ")
        respStr, _ := reader.ReadString('\n')

		//add message types
		idx := strings.IndexByte(respStr,' ')
		msgTyp := "text/plain"
		bodyStr := "";
		cont := true
		if idx>-1 {
			typ := string(respStr[:idx])
			bodyStr = string(respStr[idx+1:])
			switch typ {
				case "json":
					msgTyp = "application/json"
				case "txt":
				default:
					fmt.Printf("invalid msg type: %s\n", typ)
					cont = false
			}
		}

		cliStr = string(cliStr[:len(cliStr) -1])
		url:= dest + cliStr

 	   // Create a new request using http
		if dbg {
			fmt.Printf("url: %s\n", cliStr)
			fmt.Printf("method: %s\n", methStr)
			fmt.Printf("msgtyp: %s", msgTyp)
			fmt.Printf("body: %s", bodyStr)
		}

		if !cont {continue}

// test
//		bodyReader := bytes.NewReader([]byte(bodyStr[:len(bodyStr)]))
		reqBody := strings.NewReader(respStr)

//		if dbg {fmt.Printf("method: %s url: %s\n%s\n", methStr, url,bodyStr)}
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
			resp, err = client.Post(url, msgTyp, reqBody)
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
			req, err := http.NewRequest(methStr, url, reqBody)
			if err != nil {
				log.Printf("error -- New Req: %v\n", err)
		//        fmt.Fprintf(con, bodyStr + "\n")
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
