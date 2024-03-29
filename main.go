package main

import (
	"apiCurrency/utils"
	"fmt"
	"log"
)

func main() {
	var (
		avgCurrency, counter, min, max                 float64
		minDate, maxDate, valuteNameMin, valuteNameMax string
		err                                            error
		result                                         utils.ValCurs
	)
	urls := utils.GetDates()
	if result, err = utils.Get(urls[0]); err != nil {
		log.Printf("Failed to get XML: %v", err)
	}
	min, err = utils.ParseFloat(result.Valute[0].Value)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(urls); i++ {
		if result, err = utils.Get(urls[i]); err != nil {
			log.Printf("Failed to get XML: %v", err)
		}
		for j := 0; j < len(result.Valute); j++ {
			currentElement, err := utils.ParseFloat(result.Valute[j].Value)
			if err != nil {
				panic(err)
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
