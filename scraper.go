package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

type ShopProduct struct {
	url, image, name, price string
}

func main() {
	fmt.Println("Start scraper...")

	var shopProducts []ShopProduct
	var pagesToScrape []string

	for i := 1; i <= 48; i++ {
		pagesToScrape = append(pagesToScrape, "https://scrapeme.live/shop/page/"+strconv.Itoa(i)+"/")
	}

	c := colly.NewCollector(colly.Async(true))
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; rv:113.0) Gecko/20100101 Firefox/113.0"
	c.Limit(&colly.LimitRule{Parallelism: 4})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited: ", r.Request.URL)
	})

	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		shopProduct := ShopProduct{}

		shopProduct.url = e.ChildAttr("a", "href")
		shopProduct.image = e.ChildAttr("img", "src")
		shopProduct.name = e.ChildText("h2")
		shopProduct.price = e.ChildText(".price")

		shopProducts = append(shopProducts, shopProduct)
	})

	for _, pageToScrape := range pagesToScrape {
		c.Visit(pageToScrape)
	}

	c.Wait()

	file, err := os.Create("export/products.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"url",
		"image",
		"name",
		"price",
	}

	writer.Write(headers)

	for _, shopProduct := range shopProducts {
		record := []string{
			shopProduct.url,
			shopProduct.image,
			shopProduct.name,
			shopProduct.price,
		}
		writer.Write(record)
	}

	defer writer.Flush()

	fmt.Printf("%v", shopProducts)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
