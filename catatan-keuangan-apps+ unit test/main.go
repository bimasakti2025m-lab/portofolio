package main

import "enigmacamp.com/livecode-catatan-keuangan/delivery"

func main() {
	println("Livecode Catatan Keuangan")
	delivery.NewServer().Run()
}
