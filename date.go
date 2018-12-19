package main

import(
	"fmt"
	"time"
)

func main() {
	year, month, _ := time.Now().Date()
	t := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)

	for day := 1; day <= t.Day(); day++ {
		d := time.Date(year, month+1, day, 0, 0, 0, 0, time.UTC)
		fmt.Println(day)
		fmt.Println(d.Weekday())
	}
}
