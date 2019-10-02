package main

import (
	"fmt"
	"os"
	"log"

	"createSchema"
	"clearSchema"
)

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "create" {
			err := createSchema.CreateTables()
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("It started working again")
			}
			return
		} else if os.Args[1] == "clear" {
			err := clearSchema.ClearTables()
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Clearing database complete")
			}
			return
		}
	}
	
	fmt.Println("Nothing to do yet")
}
