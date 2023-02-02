package main

import (
	"fyne.io/fyne/v2/app"
	"training.go/Grypt/display"
)

func main() {

	// existantSecret, err := crypt.Question()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if existantSecret == 1 {

	// }
	a := app.New()
	display.Run(a)

}
