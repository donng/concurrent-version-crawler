package fetcher

import (
	"net/http"
	"golang.org/x/text/transform"
	"io/ioutil"
	"bufio"
	"golang.org/x/text/encoding/unicode"
	"log"
	"golang.org/x/text/encoding"
	"golang.org/x/net/html/charset"
	"time"
)

var rateLimiter = time.Tick(10 * time.Millisecond)

func Fetch(url string) ([]byte, error) {
	<- rateLimiter
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 判断响应编码
	newReader := bufio.NewReader(resp.Body)
	encoding := determineEncoding(newReader)
	// 转换编码格式
	utf8Reader := transform.NewReader(newReader, encoding.NewDecoder())

	// 判断返回状态码
	if resp.StatusCode != http.StatusOK {
		log.Printf("Fetcher: error response %s", resp.StatusCode)
	}

	// resp.Body 实现了 io.Reader 和 io.Closer interface
	return ioutil.ReadAll(utf8Reader)
}

// 判断内容的编码
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	// 使用缓存区，避免影响原内容的读取
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %s", err)
		return unicode.UTF8
	}

	encoding, _, _ := charset.DetermineEncoding(bytes, "")

	return encoding
}
