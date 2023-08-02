package utils

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/text/encoding/ianaindex"
)

type ValCurs struct {
	Date   string `xml:"Date,attr"`
	Name   string `xml:"name,attr"`
	Valute []struct {
		Name  string `xml:"Name"`
		Value string `xml:"Value"`
	} `xml:"Valute"`
}

func Encode(b []byte) (ValCurs, error) {
	// Декодировка в UTF-8
	var valcurs ValCurs
	decoder := xml.NewDecoder(bytes.NewBuffer(b))
	decoder.CharsetReader = func(charset string, reader io.Reader) (io.Reader, error) {
		enc, err := ianaindex.IANA.Encoding(charset)
		if err != nil {
			return nil, fmt.Errorf("charset %s: %s", charset, err.Error())
		}
		if enc == nil {
			return reader, nil
		}
		return enc.NewDecoder().Reader(reader), nil
	}
	if err := decoder.Decode(&valcurs); err != nil {
		return ValCurs{}, err
	}
	return valcurs, nil
}

func Get(url string) (ValCurs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ValCurs{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ValCurs{}, fmt.Errorf("status error: %v", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	result, err := Encode(data)
	if err != nil {
		return ValCurs{}, err
	}
	if EmptyRes(result) {
		return result, fmt.Errorf("no response for this url %s", url)
	}
	return result, nil
}
