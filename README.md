# about kaca 
a pub/sub messaging lib based on websocket  

## Getting started
```bash
go get github.com/scottkiss/kaca
```

## server

```golang
package main

import (
        "github.com/scottkiss/kaca"
       )

func main() {
    go kaca.ServeWs(":8080")
}
```

## client

```golang
package main

import (
        "fmt"
        "github.com/scottkiss/kaca"
        "time"
       )

func main() {
    producer := kaca.NewClient(":8080", "ws")
    consumer := kaca.NewClient(":8080", "ws")
    consumer.Sub("say")
    consumer.Sub("you")
    consumer.ConsumeMessage(func(message string) {
                      fmt.Println("consume =>" + message)
    })
    time.Sleep(time.Second * time.Duration(2))
    producer.Pub("you", "world")
    producer.Pub("say", "hello")
    time.Sleep(time.Second * time.Duration(2))
}


```


