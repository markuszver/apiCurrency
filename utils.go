package main

import (
	"strconv"
	"strings"
	"time"
)

func emptyRes(curs ValCurs) bool {
	if curs.Date == "" && curs.Name == "" && len(curs.Valute) == 0 {
		return true
	}
	return false
}

func parseFloat(ctringa string) (float64, error) {
	var f float64
	ctringa = strings.Replace(ctringa, ",", ".", -1)
	ctringa = strings.TrimSuffix(ctringa, "\n")
	f, err := strconv.ParseFloat(ctringa, 64)
	if err != nil {
		panic(err)
	}
	return f, nil
}

func getDates() []string {
	var result []string
	for i := 89; i >= 0; i-- {
		currentTime := time.Now()
		currentTime = currentTime.AddDate(0, 0, -i)
		url := "https://www.cbr.ru/scripts/XML_daily_eng.asp?date_req=" + currentTime.Format("02-01-2006")
		result = append(result, url)
	}
	return result
}
