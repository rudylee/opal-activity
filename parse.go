package main

import(
	"bufio"
	"encoding/csv"
	"fmt"
	"time"
	"os"
	"io"
	"log"
	"strconv"
	"strings"
)

func main() {
	// 2017
	for month := 7; month <= 12; month++ {
		parse(2017, month)
	}

	// 2018
	for month := 1; month <= 6; month++ {
		parse(2018, month)
	}
}

func parse(year int, month int) {
	t := time.Date(year, time.Month(month + 1), 0, 0, 0, 0, 0, time.UTC)

	fmt.Println(t.Month())
	csvFileName := fmt.Sprintf("export-%d-%d.csv", month - 1, year)
	csvFile, _ := os.Open(csvFileName)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	workTrips := make(map[int]bool)

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		// Convert unix timestamp to time.Time
		i, _ := strconv.ParseInt(line[1], 10, 64)
		dateTime := time.Unix(i, 0)
		_, _, day := dateTime.Date()


		if strings.Contains(line[2], "Wahroonga") ||
			 strings.Contains(line[2], "Mascot") ||
			 strings.Contains(line[2], "Hornsby") {
			workTrips[day] = true
		}
	}


	totalDaysInTheOffice := 0
	totalDaysAtHome := 0
	for day := 1; day <= t.Day(); day++ {
		d := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		weekDay := d.Weekday()

		if int(weekDay) != 6 && int(weekDay) != 0 {
			if workTrips[day] {
				totalDaysInTheOffice += 1
			} else {
				totalDaysAtHome += 1
			}
		}
	}

	fmt.Println("Total days in the office:", totalDaysInTheOffice)
	fmt.Println("Total days at home:", totalDaysAtHome)
	fmt.Println("")
}
