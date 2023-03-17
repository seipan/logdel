# logdel
This oss is a tool to delete log.Println(), etc. written for debugging, etc. 

## Example

Before
```go
package a

import (
	"log"
)

func f() {
	log.Println("nocheck") // nocheck:thislog
	log.Println("delete")
	log.Print("delete")
	if 1 == 10 {
		log.Println(11)
	}

	for i := 0; i < 2; i++ {
		log.Println(1)
	}
}
```

After
```go

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


```


## Install
```go
import "github.com/seipan/logdel/cmd/logdel"
```

## Use
```go
go vet -vettool=`which logdel` pkgname
```



