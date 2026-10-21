package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Main function called")
	fileServer := http.FileServer(http.Dir("./static"))

}
