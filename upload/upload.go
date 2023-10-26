package main

import (
	"fmt"

	"download"
)

func main() {
	fmt.Println("In upload")

	download.DownloadFile("mtgPrices.csv", "https://mtgjson.com/api/v5/csv/cardPrices.csv")
}
