package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	assetIds := []string{"1", "1027", "5426"}
	currencies := []string{"AUD"}
	q := url.Values{}
	q.Add("id", strings.Join(assetIds, ","))
	q.Add("convert", strings.Join(currencies, ","))

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "b25a7f35-2400-4192-84d7-71dba52a2cdd")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	fmt.Println(string(respBody))
}
