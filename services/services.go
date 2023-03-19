package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func getTitle(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	strBody := string(body)
	startTitle := strings.Index(strBody, "<title>")
	if startTitle == -1 {
		fmt.Println("There is no title")
		return
	}
	startTitle += 7
	endTitle := strings.Index(strBody, "</title>")
	pageTitle := []byte(strBody[startTitle:endTitle])
	fmt.Printf("Page title: %s\n", pageTitle)
}

func getWordCount(url, word string) {
	word = strings.ToLower(word)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, er := io.ReadAll(resp.Body)
	if er != nil {
		log.Fatal(er)
	}
	data := strings.ToLower(string(body))
	count := strings.Count(data, word)
	fmt.Printf("%s - %s: %d\n", url, word, count)
}

func changeWord(str string) string {
	letterDict := map[rune]string{
		'а': "a",
		'б': "b",
		'в': "v",
		'г': "g",
		'д': "d",
		'е': "e",
		'ё': "yo",
		'ж': "zh",
		'з': "z",
		'и': "i",
		'й': "y",
		'к': "k",
		'л': "l",
		'м': "m",
		'н': "n",
		'о': "o",
		'п': "p",
		'р': "r",
		'с': "s",
		'т': "t",
		'у': "u",
		'ф': "f",
		'х': "h",
		'ц': "c",
		'ч': "ch",
		'ш': "sh",
		'щ': "sh",
		'ы': "i",
		'э': "e",
		'ю': "yu",
		'я': "ya",
		' ': "-",
		'ъ': "",
		'ь': "",
	}
	str = strings.ToLower(str)

	newStr := ""

	for _, v := range str {
		newStr += letterDict[v]
	}
	return newStr
}

func getAttr(body, attr string) string {
	var attrStart int
	var attrEnd int

	switch attr {
	case "height":
		attrStart = strings.Index(body, "см")
		attrStart -= 4
		newBody := body[attrStart:]
		attrEnd = strings.Index(newBody, "\"")
		return string(newBody[:attrEnd])
	case "age":
		attrStart = strings.Index(body, "году")
		attrStart += 57
		newBody := body[attrStart:]
		attrEnd = strings.Index(newBody, ")")
		return string(newBody[:attrEnd])
	case "work":
		work := ""
		attrStart = strings.Index(body, "idstag")
		body = body[attrStart:]
		for strings.Index(body, "idstag") != -1 {
			attrStart = strings.Index(body, "idstag")
			attrStart = strings.Index(body[attrStart:], ">") + 1
			body = body[attrStart:]
			attrEnd = strings.Index(body, "<")
			work += string(body[:attrEnd]) + " "
			body = body[attrEnd+16:]
		}

		return work
	case "zodiac":
		attrStart = strings.Index(body, "{\"gen_type\":\"zodiac\"}")
		attrStart += 23
		newBody := body[attrStart:]
		attrEnd = strings.Index(newBody, "<")
		return string(newBody[:attrEnd])
	default:
		attrStart = strings.Index(body, attr)
		attrStart += len(attr) + 2
		newBody := body[attrStart:]
		attrEnd = strings.Index(newBody, "<")
		return string(newBody[:attrEnd])
	}

}

func getInfo(name string) {
	name = changeWord(name)
	url := fmt.Sprintf("https://uznayvse.ru/znamenitosti/biografiya-%s.html", name)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	a := strings.Index(string(body), "http://schema.org/Person")
	b := strings.Index(string(body), "action_block main_shadow bo_ra tablecontents")
	newBody := string(body)[a:b]
	nick := "Псевдоним: " + getAttr(newBody, "name")
	realName := getAttr(newBody, "givenName")
	patrynomic := getAttr(newBody, "additionalName")
	fullname := "Полное имя: " + realName + " " + patrynomic
	weight := "Вес: " + getAttr(newBody, "weight")
	height := "Рост: " + getAttr(newBody, "height")
	city := "Город: " + getAttr(newBody, "birthPlace")
	age := "Возраст: " + getAttr(newBody, "age")
	work := "Вид деятильности: " + getAttr(newBody, "work")
	zodiac := "Знак зодиака: " + getAttr(newBody, "zodiac")

	// fmt.Println(nick, fullname, weight, height, city, age, work, zodiac)
	wrapper(nick, fullname, height, weight, city, age, work, zodiac)
}

func wrapper(data ...interface{}) {

	fmt.Println(strings.Repeat("#", 60))
	fmt.Printf("#\n")
	for _, d := range data {

		padding := 10

		fmt.Printf("#%s%v\n", strings.Repeat(" ", padding), d)
		fmt.Printf("#\n")

	}

	fmt.Println(strings.Repeat("#", 60))
	fmt.Println()
	fmt.Println()
}
