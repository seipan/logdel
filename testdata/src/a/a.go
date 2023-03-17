package a

import (
	"log"
)

func f() {
	log.Println("nocheck")
	if 1 == 10 {
	}
	for i := 0; i < 2; i++ {
	}
}

// this func is ......
//aa
// nocheck:thislog
