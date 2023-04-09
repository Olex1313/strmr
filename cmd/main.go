package main

import "gobs/internal"

func main() {
	c := internal.NewClient()
	internal.NewServer(c.GetStream)
}
