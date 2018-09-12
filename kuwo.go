package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func main() {

	id := os.Args[1]

	fmt.Println(id)

	url := fmt.Sprintf("http://www.kuwo.cn/artist/contentMusicsAjax?artistId=%s&pn=1&rn=1000", id)

	fmt.Println(url)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	html, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Printf("%s", html)

	exp := regexp.MustCompile(`data-music='(.*)'>`)

	rs := exp.FindAllStringSubmatch(string(html), -1) //[3.14 1.0 6.66]

	for _, item := range rs {
		fmt.Println(item[1])

		data := map[string]string{}

		err := json.Unmarshal([]byte(item[1]), &data)

		if err != nil {
			fmt.Println(err)
			continue
		}

		download(data["id"], data["name"], data["artist"])
	}

}

func download(id string, name string, artist string) {

	fileName := fmt.Sprintf("%s/Downloads/%s - %s.mp3", os.Getenv("HOME"), artist, name)
	fmt.Println(fileName)

	format := "mp3"
	url := fmt.Sprintf("http://antiserver.kuwo.cn/anti.s?format=%s&rid=%s&type=convert_url&response=res", format, id)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	mp3, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	}

	f.Write(mp3)

	if err := f.Close(); err != nil {
		fmt.Println(err)
	}

}
