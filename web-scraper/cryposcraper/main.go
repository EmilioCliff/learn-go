package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"encoding/csv"
	"os"
)

type CryptocurrencyRecord struct {
	Name string
	Symbol string
	Price string
	Volume24h string
	Change1h string
	Change24h string
	Change7d string
}

func main() {
	c := colly.NewCollector()
	var records []CryptocurrencyRecord

	c.OnRequest(func(r *colly.Request){
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		var record CryptocurrencyRecord
		record.Name = e.ChildText("td.cmc-table__cell--sort-by__name > div > a.cmc-table__column-name--name")
		record.Symbol = e.ChildText("td.cmc-table__cell--sort-by__symbol > div")
		record.Price = e.ChildText("td.cmc-table__cell--sort-by__price > div > a > span")
		record.Volume24h = e.ChildText("td.cmc-table__cell--sort-by__volume-24-h > a")
		record.Change1h = e.ChildText("td.cmc-table__cell--sort-by__percent-change-1-h > div")
		record.Change24h = e.ChildText("td.cmc-table__cell--sort-by__percent-change-24-h > div")
		record.Change7d = e.ChildText("td.cmc-table__cell--sort-by__percent-change-7-d > div")

		records = append(records, record)
	})

	c.Visit("https://coinmarketcap.com/all/views/all/")
	file, err := os.Create("output.csv")
	if err != nil {
		fmt.Println("Error :", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, r := range records{
		new := []string{
			r.Name,
			r.Symbol,
			r.Price,
			r.Volume24h,
			r.Change1h,
			r.Change24h,
			r.Change7d,
		}

		writer.Write(new)
	}
}