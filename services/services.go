package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func GetTitle(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	strBody := string(body)
	startTitle := strings.Index(strBody, "<title>")
	//Check if startTitle equals -1 then there is no title
	if startTitle == -1 {
		fmt.Println("There is no title")
		return
	}
	startTitle += 7
	endTitle := strings.Index(strBody, "</title>")
	pageTitle := []byte(strBody[startTitle:endTitle])
	fmt.Printf("Page title: %s\n", pageTitle)
}

func GetWordCount(url, word string) {
	//Make all words lowercase
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
	//Converts russian letters into english ones
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

		// return strings.Replace(strings.Trim(work, " "), " ", ", ", -1)
		return work
	case "zodiac":
		attrStart = strings.Index(body, "{\"gen_type\":\"zodiac\"}")
		attrStart += 23
		newBody := body[attrStart:]
		attrEnd = strings.Index(newBody, "<")
		return string(newBody[:attrEnd])
	case "likes":
		attrStart = strings.Index(body, "data-positive")
		attrStart += 15
		newBody := body[attrStart:]
		attrEnd = strings.Index(newBody, "\"")
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", 92, newBody[:attrEnd])
	case "dislikes":
		attrStart = strings.Index(body, "data-negative")
		attrStart += 15
		newBody := body[attrStart:]
		attrEnd = strings.Index(newBody, "\"")
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", 91, newBody[:attrEnd])
	default:
		attrStart = strings.Index(body, attr)
		//checking if patrynomic is none
		if attrStart == -1 {
			return "-"
		}
		attrStart += len(attr) + 2
		newBody := body[attrStart:]
		attrEnd = strings.Index(newBody, "<")
		return string(newBody[:attrEnd])
	}

}

func GetInfo(name string) {
	name = changeWord(name)
	url := fmt.Sprintf("https://uznayvse.ru/znamenitosti/biografiya-%s.html", name)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	} else if resp.StatusCode == 404 {
		fmt.Printf("%s не был найден\n", name)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	startBody := strings.Index(string(body), "http://schema.org/Person")
	endBody := strings.Index(string(body), "action_block main_shadow bo_ra tablecontents")
	newBody := string(body)[startBody:endBody]
	nick := getAttr(newBody, "name")
	// realName := getAttr(newBody, "givenName")
	patrynomic := getAttr(newBody, "additionalName")
	fullname := "Полное имя: " + nick + " " + patrynomic
	weight := "Вес: " + getAttr(newBody, "weight")
	height := "Рост: " + getAttr(newBody, "height")
	city := "Город: " + getAttr(newBody, "birthPlace")
	age := "Возраст: " + getAttr(newBody, "age")
	work := "Вид деятильности: " + getAttr(newBody, "work")
	zodiac := "Знак зодиака: " + getAttr(newBody, "zodiac")
	likes := "Нравится: " + getAttr(newBody, "likes")
	dislikes := "Не нравится: " + getAttr(newBody, "dislikes")

	sex := ""
	//checking if last letter in patrynomic equals russian "а" letter
	if patrynomic[len(patrynomic)-1] != 176 {
		sex = "boy"
	} else {
		sex = "girl"
	}

	//Call the function which obtains and returns formatted data
	wrapper(fullname, height, weight, city, age, work, zodiac, likes, dislikes, sex)
}

func wrapper(data ...interface{}) {
	var color uint8
	if data[len(data)-1] == "boy" {
		//stands for blue
		color = 94
	} else {
		//stands for pink
		color = 95
	}
	//returns colored "#" symbol which depends on sex
	coloredSharp := fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, "#")
	fmt.Println()
	fmt.Println(strings.Repeat(coloredSharp, 60))
	fmt.Printf("%s\n", coloredSharp)

	for i, d := range data {
		if i == len(data)-1 {
			break
		}
		//Padding from left side
		padding := 5
		fmt.Printf("%s%s%v\n", coloredSharp, strings.Repeat(" ", padding), d)
		fmt.Printf("%s\n", coloredSharp)
	}

	fmt.Println(strings.Repeat(coloredSharp, 60))
	fmt.Println()
	fmt.Println()
}
