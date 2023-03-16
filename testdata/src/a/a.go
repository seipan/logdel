package a

import (
	"log"
)

func f() {
	log.Println("hi") //aa
	log.Println("hi") // nocheck:thislog
	log.Print("hi")
	if 1 == 10 {
		log.Println(11)
	}

	for i := 0; i < 2; i++ {
		log.Println(1)
	}

}
