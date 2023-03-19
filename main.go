package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"parser/services/services"
	"time"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Print("Кого ищите?: ")
	name, err := in.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	t := time.Now()
	services.GetInfo(name)

	fmt.Printf("Время работы: %s\n", time.Since(t))
}
