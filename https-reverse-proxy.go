// A TLS-terminating single-backend reverse proxy. Listens on the address given
// with the --from flag and forwards all traffic to the server given with the
// --to flag.
// Similar to basic-reverse-proxy, but talks HTTPS.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.

// Borrowed and added to.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Listen_Address string   `yaml:"listen"`
	Server         string   `yaml:"server"`
	Cert_File      string   `yaml:"cert_file"`
	Key_File       string   `yaml:"key_file"`
	Paths          []string `yaml:"paths"`
	Default_Domain string   `yaml:"default_domain"`
}

func main() {
	file := flag.String("f", "", "Specify config yaml file")

	flag.Parse()

	if !isFlagPassed("f") {
		usage()
		os.Exit(1)
	}

	config := ChkYaml(file)

	listen_address := config.Listen_Address
	server := config.Server
	cert_file := config.Cert_File
	key_file := config.Key_File

	paths := config.Paths
	default_domain := config.Default_Domain

	log.Println("Default domain: ", default_domain)

	final_url := ParseToUrl(server)
	proxy_server := httputil.NewSingleHostReverseProxy(final_url)

	regular_web := ParseToUrl(default_domain)
	proxy_web := httputil.NewSingleHostReverseProxy(regular_web)

	for _, path := range paths {
		log.Println("Path: ", path)
		handle_path := path
		http.Handle(handle_path, proxy_server)
	}

	http.Handle("/", proxy_web) // Varies with self signed certificates for testing and prolly all when it comes to cloudflare/akamai/other-bollocks

	log.Println("Starting proxy server on", listen_address)
	err := http.ListenAndServeTLS(listen_address, cert_file, key_file, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s: \n", os.Args[0])
	fmt.Println()
	flag.PrintDefaults()
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func ChkYaml(file *string) Config {
	var config Config
	_, err := os.Stat(*file)
	if err == nil {
		data, err := os.ReadFile(*file)
		if err != nil {
			panic(err)
		}

		if err := yaml.Unmarshal(data, &config); err != nil {
			panic(err)
		}
	}
	return config

}

// ParseToUrl parses a "to" address to url.URL value
func ParseToUrl(addr string) *url.URL {
	if !strings.HasPrefix(addr, "https") {
		addr = "https://" + addr
	}
	final_url, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}
	return final_url
}
