package main

import (
	"apiCurrency/utils"
	"fmt"
	"log"
	"net/http"
)

type result struct {
	avgCurrency, min, max                          float64
	minDate, maxDate, valuteNameMin, valuteNameMax string
}

func main() {
	var res result
	ch := make(chan struct{})
	go calculation(ch, 90, &res)
	http.HandleFunc("/apiCurrency", func(w http.ResponseWriter, req *http.Request) {
		select {
		case <-ch:
			fmt.Fprintf(w, "Максимальное значение %f, название %s, дата %s\n", res.max, res.valuteNameMax, res.maxDate)
			fmt.Fprintf(w, "Минимальное значение %f, название %s, дата %s\n", res.min, res.valuteNameMin, res.minDate)
			fmt.Fprintf(w, "Среднее значение курса рубля за весь период по всем валютам: %f", res.avgCurrency)
			return
		default:
			fmt.Fprint(w, "Идет обработка, подождите :)\n")
		}
	})
	http.ListenAndServe(":8080", nil)
}

func calculation(ch chan struct{}, days int, res *result) {
	defer close(ch)
	var (
		avgCurrency, counter, min, max                 float64
		minDate, maxDate, valuteNameMin, valuteNameMax string
		err                                            error
		response                                       utils.ValCurs
	)
	urls := utils.GetDates(days)
	if response, err = utils.Get(urls[0]); err != nil {
		log.Printf("Failed to get XML: %v", err)
	}
	min, err = utils.ParseFloat(response.Valute[0].Value)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(urls); i++ {
		if response, err = utils.Get(urls[i]); err != nil {
			log.Printf("Failed to get XML: %v", err)
		}
		for j := 0; j < len(response.Valute); j++ {
			currentElement, err := utils.ParseFloat(response.Valute[j].Value)
			if err != nil {
				panic(err)
			}
			if min > currentElement {
				min = currentElement
				valuteNameMin = response.Valute[j].Name
				minDate = response.Date
			}
			if max < currentElement {
				max = currentElement
				valuteNameMax = response.Valute[j].Name
				maxDate = response.Date
			}
			avgCurrency += currentElement
			counter++
		}
	}
	avgCurrency /= counter
	res.avgCurrency = avgCurrency
	res.max = max
	res.minDate = minDate
	res.maxDate = maxDate
	res.min = min
	res.valuteNameMax = valuteNameMax
	res.valuteNameMin = valuteNameMin
	fmt.Print("channel was closed")
}
