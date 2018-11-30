package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	username := flag.String("username", "", "Your opal.com.au username")
	password := flag.String("password", "", "Your opal.com.au password")
	month := flag.Int("month", -1, "The month of activity")
	year := flag.Int("year", -1, "The year of activity")

	flag.Parse()

	if *username == "" {
		log.Fatal("Please provide username")
	}

	if *password == "" {
		log.Fatal("Please provide password")
	}

	if *month == -1 {
		log.Fatal("Please provide month")
	}

	if *year == -1 {
		log.Fatal("Please provide year")
	}

	Login("cookieJar", *username, *password)

	GetMonthlyActivity(*month, *year, 1)
}

func GetMonthlyActivity(month int, year int, page int) {
	fmt.Printf("Fetching Month: %d Year: %d Page: %d\n", month, year, page)

	activityUrl := fmt.Sprintf("https://www.opal.com.au/registered/opal-card-transactions/opal-card-activities-list?AMonth=%d&AYear=%d&cardIndex=0&pageIndex=%d", month, year, page)

	d, _ := http.NewRequest("GET", activityUrl, nil)

	client := &http.Client{
		Jar: CookieJar(activityUrl),
	}

	resp, _ := client.Do(d)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	csvFileName := fmt.Sprintf("export-%d-%d.csv", month, year)
	csvFile, _ := os.OpenFile(csvFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Parse the activity
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		/*
			Column 0 - Transaction Number
			Column 1 - Date/time
			Column 2 - Mode ( Ignored )
			Column 3 - Details ( From and to )
			Column 4 - Journey Number
			Column 5 - Fare Applied eg: Travel Reward, Day Cap
			Column 6 - Discount
			Column 7 - Amount
		*/
		csvRow := []string{}
		s.Find("td").Each(func(j int, row *goquery.Selection) {
			if j == 0 {
				csvRow = append(csvRow, strings.TrimSpace(row.Text()))
			} else if j == 1 {
				travelDateTime := strings.TrimSpace(row.Text())

				csvRow = append(csvRow, ConvertTravelDateTimeToUnix(travelDateTime))
			} else if j != 2 {
				csvRow = append(csvRow, strings.TrimSpace(row.Text()))
			}
		})
		writer.Write(csvRow)
	})

	// Find the pagination
	nextPage, found := doc.Find("a#next").Attr("href")

	if found {
		// Sleep for 5 second so we are not hammering Opal website
		time.Sleep(5 * time.Second)
		pageIndexRegex := regexp.MustCompile(`pageIndex=(.*)`)
		match := pageIndexRegex.FindStringSubmatch(nextPage)

		page, _ = strconv.Atoi(match[1])
		GetMonthlyActivity(month, year, page)
	}
}

func ConvertTravelDateTimeToUnix(dateString string) string {
	layout := "Mon02/01/200615:04 MST"
	t, _ := time.Parse(layout, dateString+" AEST")
	return strconv.FormatInt(t.Unix(), 10)
}

func Login(cookieFile string, username string, password string) {
	fmt.Printf("Logging in as %s\n", username)
	file, err := os.Create(cookieFile)

	if err != nil {
		log.Fatal("Unable to create cookie file")
	}

	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}

	resp, _ := client.PostForm("https://www.opal.com.au/login/registeredUserUsernameAndPasswordLogin", url.Values{
		"h_username": {username},
		"h_password": {password},
	})
	resp.Body.Close()

	for _, cookie := range cookieJar.Cookies(resp.Request.URL) {
		if cookie.Name == "JSESSIONID" {
			file.WriteString(cookie.Value)
		}
	}
}

func CookieJar(targetUrl string) *cookiejar.Jar {
	file, err := os.Open("cookieJar")
	if err != nil {
		log.Fatal("Unable to open cookie jar file")
	}

	cookieData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	cookieJar, _ := cookiejar.New(nil)

	var cookies []*http.Cookie

	session_cookie := &http.Cookie{
		Name:   "JSESSIONID",
		Value:  string(cookieData),
		Path:   "/",
		Domain: "www.opal.com.au",
	}

	cookies = append(cookies, session_cookie)

	cookieURL, _ := url.Parse(targetUrl)

	cookieJar.SetCookies(cookieURL, cookies)

	return cookieJar
}
