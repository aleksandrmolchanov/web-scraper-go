package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type ShopProduct struct {
	url, image, name, price string
}

func main() {
	fmt.Println("Start scraper...")

	var shopProducts []ShopProduct

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	var nodes []*cdp.Node
	chromedp.Run(ctx,
		chromedp.Navigate("https://scrapeme.live/shop/"),
		chromedp.Nodes(".product", &nodes, chromedp.ByQueryAll),
	)

	var url, image, name, price string
	for _, node := range nodes {
		chromedp.Run(ctx,
			chromedp.AttributeValue("a", "href", &url, nil, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.AttributeValue("img", "src", &image, nil, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text("h2", &name, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(".price", &price, chromedp.ByQuery, chromedp.FromNode(node)),
		)

		shopProduct := ShopProduct{}

		shopProduct.url = url
		shopProduct.image = image
		shopProduct.name = name
		shopProduct.price = price

		shopProducts = append(shopProducts, shopProduct)
	}

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
