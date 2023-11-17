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
	var pagesToScrape []string

	pageToScrape := "https://scrapeme.live/shop/page/1/" 
	pagesDiscovered := []string{ pageToScrape }
	i := 1
	limit := 5

	c := colly.NewCollector()

	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		newPaginationLink := e.Attr("href")

		/*if !contains(pagesToScrape, newPaginationLink) {
			if !contains(pagesDiscovered, newPaginationLink) {
				pagesToScrape = append(pagesToScrape, newPaginationLink)
			}
			pagesDiscovered = append(pagesDiscovered, newPaginationLink)
		}*/

		if !contains(pagesDiscovered, newPaginationLink) {
			pagesToScrape = append(pagesToScrape, newPaginationLink)
			pagesDiscovered = append(pagesDiscovered, newPaginationLink)
		}
			
	})

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

	c.OnScraped(func(response *colly.Response) {
		if len(pagesToScrape) != 0 && i < limit {
			pageToScrape = pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]

			i++

			c.Visit(pageToScrape)
		}
	})
	
	c.Visit(pageToScrape)

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