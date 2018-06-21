package main

func main() {
	a := App{}

	// Set Username and Password here
	a.Initialize("pentre", "123", "")

	a.Run(":8080")
}
