/*
Утилита wget

Реализовать утилиту wget с возможностью скачивать сайты целиком.
*/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	url "net/url"
)

// launch_wget создает рабочую директорию и выполняет wget
func launch_wget(url string, outPath string, depth int, timeout time.Duration) {
	if depth < 1 {
		log.Fatal("Некорректное значение")
	}

	wd, _ := os.Getwd()
	err := os.Mkdir(wd+`/`+outPath, os.ModeDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal(err)
		}
	}
	err = os.Chdir(outPath)
	if err != nil {
		log.Fatal("Невозможно изменить рабочий каталог")
	}
	make_recursion(url, url, depth, 0, 1, timeout)
}

// loadPage загружает и сохраняет страницу по url
func loadPage(address string) []byte {
	resp, err := http.Get(address)
	if err != nil {
		log.Println("Невозможно загрузить страницу по URL", address)
		return nil
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Невозможно прочитать данные страницы с URL-адресом", address)
	}

	return data
}

// writeInformation записывает информацию в файл
func writeInformation(data []byte, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// make_recursion реализет wget
func make_recursion(rootUrl string, curUrl string, depth int, curDepth int, pageInd int, timeout time.Duration) {
	if curDepth == depth {
		return
	}

	data := loadPage(curUrl)

	if data != nil {
		links := sliceOfLinks(data)

		linksNorm := make([]string, len(links))
		copy(linksNorm, links)
		formatLinks(linksNorm, rootUrl, curUrl)

		for i := 0; i < len(links); i++ {
			oldLink := []byte(links[i])
			newLink := []byte(strconv.Itoa(curDepth+1) + "_" + strconv.Itoa(i) + ".html")
			copy(data, bytes.ReplaceAll(data, oldLink, newLink))
		}

		err := writeInformation(data, strconv.Itoa(curDepth)+"_"+strconv.Itoa(pageInd)+".html")
		if err != nil {
			log.Println("Невозможно записать файл: ", err)
		}

		for i := 0; i < len(linksNorm); i++ {
			time.Sleep(timeout)
			make_recursion(rootUrl, linksNorm[i], depth, curDepth+1, i, timeout)
		}
	}
}

// sliceOfLinks возвращает слайс ссылок из тела заданной страницы
func sliceOfLinks(body []byte) []string {
	var links []string
	bodyReader := bytes.NewReader(body)
	z := html.NewTokenizer(bodyReader)
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}

				}
			}

		}
	}
}

// AddSuffix добавляет / в конец ссылки
func AddSuffix(url string) string {
	if !strings.HasSuffix(url, "/") {
		return url + "/"
	} else {
		return url
	}
}

// formatLinks приводит ссылки к полному виду (с протоколом и доменом)
func formatLinks(links []string, rootUrl string, parentUrl string) {
	for i, link := range links {
		if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
			if strings.HasPrefix(link, "/") {
				url, err := url.Parse(rootUrl)
				if err != nil {
					log.Println("Невозможно разобрать URL: ", url)
				}
				links[i] = url.Scheme + "://" + AddSuffix(url.Host) + link[1:]
			} else {
				links[i] = AddSuffix(parentUrl) + link
			}
		}
	}
}

func main() {
	scan := bufio.NewScanner(os.Stdin)

	fmt.Println("Введите полный URL")
	scan.Scan()
	url := scan.Text()

	fmt.Println("Введите полный путь к каталогу для хранения html-файлов")
	scan.Scan()
	path := scan.Text()

	launch_wget(url, path, 2, 3)

	fmt.Println("Загрузка завершена")
}
