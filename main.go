package main

import (
	"fmt"
	"log"
)

func main() {
	var avgCurrency, counter float64
	var min, max float64
	var minDate, maxDate, valuteNameMin, valuteNameMax string
	var err error
	var result ValCurs
	var urls []string
	urls = getDates()
	for i := 0; i < len(urls); i++ {
		if result, err = get(urls[i]); err != nil {
			log.Printf("Failed to get XML: %v", err)
		}
		min, err = parseFloat(result.Valute[0].Value)
		if err != nil {
			panic(err)
			return
		}
		max, err = parseFloat(result.Valute[0].Value)
		if err != nil {
			panic(err)
			return
		}
		for j := 0; j < len(result.Valute); j++ {
			currentElement, err := parseFloat(result.Valute[j].Value)
			if err != nil {
				panic(err)
				return
			}
			if min > currentElement {
				min = currentElement
				valuteNameMin = result.Valute[j].Name
				minDate = result.Date
			}
			if max < currentElement {
				max = currentElement
				valuteNameMax = result.Valute[j].Name
				maxDate = result.Date
			}
			avgCurrency += currentElement
			counter++
		}
	}
	avgCurrency /= counter
	fmt.Printf("Максимальное значение %f, название %s, дата %s\n", max, valuteNameMax, maxDate)
	fmt.Printf("Минимальное значение %f, название %s, дата %s\n", min, valuteNameMin, minDate)
	fmt.Printf("Среднее значение курса рубля за весь период по всем валютам: %f", avgCurrency)
}
