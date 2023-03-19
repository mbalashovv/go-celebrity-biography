package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	// go getTitle("https://migrantumir.com/reyting-stran-mira-po-urovnyu-zhizni")
	// go getTitle("https://mpt.ru")
	// getTitle("http://81.200.151.200")
	// getTitle("https://citilink.ru")
	// go getWordCount("https://mpt.ru", "а")
	// go getWordCount("https://www.rbc.ru", "а")
	// getWordCount("https://kino.ru", "a")
	services.getInfo("Катя Клэп")
	fmt.Printf("Took %s\n", time.Since(t))
}
