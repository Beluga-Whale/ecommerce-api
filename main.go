package main

import "github.com/Beluga-Whale/ecommerce-api/config"

func main() {
	// NOTE - LoadEnv
	config.LoadEnv()

	// NOTE - Connect DB
	config.ConnectDB()

	// // NOTE - Fiber

	// app :=
}