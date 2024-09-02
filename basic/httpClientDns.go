// program that retrieves the source from a url
// modified from https://koraygocmen.medium.com/custom-dns-resolver-for-the-default-http-client-in-go-a1420db38a5d
//
// added cli
//
// todo add resolver as option
//

package main

import (
	"context"
	"io/ioutil"
	"os"
	"log"
	"fmt"
	"net"
	"net/http"
	"time"
	util "github.com/prr123/utility/utilLib"
)

func main() {

    numarg := len(os.Args)
    flags:=[]string{"dbg", "url", "dns"}

    useStr := " /url=urlstr [/dns=dnsresolv] [/dbg]"
    helpStr := "server"

    if numarg > len(flags) +1 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: %s %s\n", os.Args[0], useStr)
        os.Exit(-1)
    }

    if numarg == 1 || (numarg > 1 && os.Args[1] == "help") {
        fmt.Printf("help: %s %s\n", os.Args[0], helpStr)
        fmt.Printf("usage is: %s %s\n", os.Args[0], useStr)
        os.Exit(1)
    }

    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}

    dbg:= false
    _, ok := flagMap["dbg"]
    if ok {dbg = true}

    urlStr := ""
    uval, ok := flagMap["url"]
    if !ok {
        log.Fatalf(" error no url provided!\n")
    } else {
        if uval.(string) == "none" {log.Fatalf("error: no url value provided!\n")}
        urlStr = uval.(string)
    }

    dnsStr := "google"
    dval, ok := flagMap["dns"]
    if ok {
        if dval.(string) != "none" {urlStr = uval.(string)}
    }

	dnsIP:=""
	switch dnsStr {
		case "google":
			dnsIP = "8.8.8.8"
		case "cloudflare":
			dnsIP = "1.1.1.1"
		default:
			log.Fatalf("unknown dns resolver: %s\n", dnsStr)
	}

	if dbg {
    	log.Printf("debug -- url: %s\n", urlStr)
    	log.Printf("debug -- dns: %s\n", dnsStr)
	}

	var (
    	dnsResolverIP        = dnsIP+":53" // Google DNS resolver.
		dnsResolverProto     = "udp"        // Protocol to use for the DNS resolver
		dnsResolverTimeoutMs = 5000         // Timeout (ms) for the DNS resolver (optional)
  	)

	dialer := &net.Dialer{
    	Resolver: &net.Resolver{
      		PreferGo: true,
      		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
        		d := net.Dialer{
          			Timeout: time.Duration(dnsResolverTimeoutMs) * time.Millisecond,
        			}
				return d.DialContext(ctx, dnsResolverProto, dnsResolverIP)
      		},
    	},
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
    	return dialer.DialContext(ctx, network, addr)
  	}

	http.DefaultTransport.(*http.Transport).DialContext = dialContext
	httpClient := &http.Client{}

  	// Testing the new HTTP client with the custom DNS resolver.
	url := "https://" + urlStr
	if dbg {log.Printf("debug -- url: %s\n", url)}
  	resp, err := httpClient.Get(url)
  	if err != nil {log.Fatalf("error -- http get: %v\n", err)}
  	defer resp.Body.Close()

  	body, err := ioutil.ReadAll(resp.Body)
  	if err != nil {
    	log.Fatalf("error -- read body: %v\n",err)
  	}

	fmt.Println("***** body *****")
	fmt.Println(string(body))
	fmt.Println("*** end body ***")

	log.Println("success")
}
