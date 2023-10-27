package main

import (
	"godis/tcp"

	RedisServer "godis/redis/server"
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
	}, RedisServer.MakeServer())

	if err != nil {
	}
}
