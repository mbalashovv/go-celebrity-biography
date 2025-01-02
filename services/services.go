package services

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

func changeWord(str string) string {
	letterDict := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d",
		'е': "e", 'ё': "yo", 'ж': "zh", 'з': "z", 'и': "i",
		'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n",
		'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t",
		'у': "u", 'ф': "f", 'х': "h", 'ц': "c", 'ч': "ch",
		'ш': "sh", 'щ': "sh", 'ы': "i", 'э': "e", 'ю': "yu",
		'я': "ya", ' ': "-", 'ъ': "", 'ь': "",
	}
	str = strings.TrimSpace(strings.ToLower(str))

	var newStr strings.Builder
	for _, v := range str {
		if value, exists := letterDict[v]; exists {
			newStr.WriteString(value)
		} else {
			newStr.WriteString(string(v))
		}
	}
	return newStr.String()
}

func GetInfo(nameToFind string) {
	// transliterate russian name to english one
	nameToFind = changeWord(nameToFind)

	url := fmt.Sprintf("https://uznayvse.ru/znamenitosti/biografiya-%s.html", nameToFind)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	} else if resp.StatusCode == 404 {
		fmt.Printf("%s не был найден\n", nameToFind)
		return
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	infoBox := doc.Find(`div[itemtype="http://schema.org/Person"]`)

	// Extract data using selectors

	// Name could be a list of elements, therefore check it for lenght
	nameElement := infoBox.Find("span[itemprop='name']")
	name := ""
	if nameElement.Length() > 1 {
		name = nameElement.First().Text()
	} else {
		name = nameElement.Text()
	}

	patrynomic := infoBox.Find("strong[itemprop='additionalName']").Text()
	profession := "Вид деятельности: " + strings.TrimSpace(infoBox.Find("dt:contains('Кто такой:')").Next().Find("strong").Text())
	birthDate := "Датя рождения: " + infoBox.Find("meta[itemprop='birthDate']").AttrOr("content", "N/A")
	ageParts := strings.Split(infoBox.Find("dt:contains('День рождения:')").Next().Text(), "(")
	age := "Возраст: " + strings.TrimSuffix(strings.TrimSpace(ageParts[len(ageParts)-1]), ")")
	city := "Город: " + infoBox.Find("strong[itemprop='birthPlace']").Text()
	height := "Рост: " + infoBox.Find("strong[itemprop='height']").Text()
	weight := "Вес: " + infoBox.Find("strong[itemprop='weight']").Text()
	zodiac := "Знак зодиака: " + infoBox.Find("dt:contains('Знак Зодиака:')").Next().Find("a").First().Text()
	likes := "Нравится: " + infoBox.Find("[data-positive]").AttrOr("data-positive", "-")
	dislikes := "Не нравится: " + infoBox.Find("[data-negative]").AttrOr("data-negative", "-")

	fullname := "Полное имя: " + name + " " + patrynomic

	sex := "male"
	//checking if the last letter in patrynomic equals russian "а" letter
	if len(patrynomic) > 0 && patrynomic[len(patrynomic)-1] == 176 {
		sex = "female"
	}

	wrapper([]interface{}{fullname, height, weight, city, age, birthDate, profession, zodiac, likes, dislikes}, sex)
}

func wrapper(data []interface{}, sex string) {
	var c *color.Color
	if sex == "male" {
		c = color.New(color.FgHiBlue)
	} else {
		c = color.New(color.FgHiMagenta)
	}

	coloredSharp := c.Sprint("#")
	fmt.Println()
	fmt.Println(strings.Repeat(coloredSharp, 60))
	fmt.Printf("%s\n", coloredSharp)

	for _, d := range data {
		padding := 5
		fmt.Printf("%s%s%v\n", coloredSharp, strings.Repeat(" ", padding), d)
		fmt.Printf("%s\n", coloredSharp)
	}

	fmt.Println(strings.Repeat(coloredSharp, 60))
	fmt.Println()
}
