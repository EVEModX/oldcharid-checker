package main

import (
        "fmt"
	"os"
	"log"	
        "io/ioutil"
        "net/http"
	"flag"
)

var concurrency *int = flag.Int("c", 10, "Default 10")
var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)

func Get(url string) (content string, statusCode int) {
        resp, err1 := http.Get(url)
        if err1 != nil {
                statusCode = -100
                return
        }
        defer resp.Body.Close()
        data, err2 := ioutil.ReadAll(resp.Body)
        if err2 != nil {
                statusCode = -200
                return
        }
        statusCode = resp.StatusCode
        content = string(data)
        return
}

func main() {

    dstFile,err := os.Create("charid.txt")
    if err!=nil{
        fmt.Println(err.Error())    
        return
    }   	
defer dstFile.Close()
	flag.Parse()
	flag.Parse()
	//concurrency := 10
	sem := make(chan bool, *concurrency)

	baseUrl := "https://api.eve-online.com.cn/eve/characterinfo.xml.aspx?characterid="

	for i := 100000000; i < 999999999; i++{
		sem <- true
		go func(i int){
			defer func() { <-sem }()
			_, code := Get(baseUrl + fmt.Sprintf("%d", i))
			var check string
			if code == 200 {
				dstFile.WriteString(fmt.Sprintf("%d", i) + "\n")
				check = "yes"
			} else {
				check = "no"
			}
			logger.Println(fmt.Sprintf("[INFO] Checking %d ... %s", i, check))
		}(i)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

