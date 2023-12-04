package main

import (
	"godis/tcp"

	"godis/redis/net"
)

var banner = `
   ______          ___
  / ____/___  ____/ (_)____
 / / __/ __ \/ __  / / ___/
/ /_/ / /_/ / /_/ / (__  )
\____/\____/\__,_/_/____/
`

func main() {
	print(banner)

	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: "0.0.0.0:6639",
	}, net.MakeHandler())

	if err != nil {
	}
}
