package a

import (
	"log"
)

func f() {
	log.Println("delete")
	log.Print("delete")
	if 1 == 10 {
		log.Println(11)
	}

	for i := 0; i < 2; i++ {
		log.Println(1)
	}
}
