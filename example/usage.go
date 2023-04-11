package main

import (
	"fmt"

	"github.com/gol-gol/golimg"
)

var txt = `
At the time, no single team member knew Go, but within a month, everyone was writing in Go and we were building out the endpoints. It was the flexibility, how easy it was to use, and the really cool concept behind Go (how Go handles native concurrency, garbage collection, and of course safety+speed.) that helped engage us during the build. Also, who can beat that cute mascot!
`

func main() {
	drawtxt := golimg.DrawText{
		FontPath: "fonts/FFF_Tusj.ttf",
		FontSize: 18.0,
	}
	saveAs := "out.png"
	drawtxt.CreateImageWithText(txt, saveAs)
	fmt.Println("check:", saveAs)
}
