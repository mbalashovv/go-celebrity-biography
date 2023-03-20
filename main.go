package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"parser/services"
)

func main() {
	var running = true
	t := time.Now()
	for running {
		in := bufio.NewReader(os.Stdin)
		fmt.Print("Поиск: ")
		name, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		services.GetInfo(name)

		var input string
		fmt.Print("Продолжить работу?(д/н): ")
		fmt.Scan(&input)
		if input == "д" {
			continue
		} else {
			running = false
		}
	}
	fmt.Printf("Время работы: %s\n", time.Since(t))

}
