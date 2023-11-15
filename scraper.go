package main

import (
	"fmt"
	"github.com/gocolly/colly" 
	"encoding/csv" 
	"log" 
	"os" 
)

type ShopProduct struct {
	url, image, name, price string
}

func main() {
	fmt.Println("Start scraper...")

	var shopProducts []ShopProduct

	c := colly.NewCollector()

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
	
	c.Visit("https://scrapeme.live/shop/")

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