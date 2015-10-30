# about kaca 
a pub/sub messaging system based on websocket  

## Getting started
```bash
go get github.com/scottkiss/kaca
```

## server

```go
package main

import (
        "github.com/scottkiss/kaca"
       )

func main() {
    //use true to set check origin
    kaca.ServeWs(":8080",true)
}
```

## pub/sub client

```go
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

## broadcast client
```go

}ckage main

import (
        "fmt"
        "github.com/scottkiss/kaca"
        "time"
       )

func main() {
              producer := kaca.NewClient(":8080", "ws")
              consumer := kaca.NewClient(":8080", "ws")
              c2 := kaca.NewClient(":8080", "ws")
              c2.ConsumeMessage(func(message string) {
                      fmt.Println("c2 consume =>" + message)
                      })
              consumer.Sub("say")
              consumer.Sub("you")
              consumer.ConsumeMessage(func(message string) {
                      fmt.Println("consume =>" + message)
                      })
              time.Sleep(time.Second * time.Duration(2))
              producer.Broadcast("broadcast...")
              time.Sleep(time.Second * time.Duration(2))
})
}

```

