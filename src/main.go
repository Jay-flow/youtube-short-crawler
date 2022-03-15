package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fedesog/webdriver"
	"github.com/tebeka/selenium"
)

var session webdriver.Session

func main() {
	videoCounts := 100

	session = setupWebDriver()

	f, w := createTheCSV()

	goToTheShort()

	for i := 1; i <= videoCounts; i++ {
		currentURL := getCurrnetShortUrl()
		fmt.Println(currentURL)

		rawURL := strings.Split(currentURL, "/")
		data := []string{rawURL[len(rawURL)-1], currentURL}

		w.Write(data)
		w.Flush()

		time.Sleep(1 * time.Second)
		goToTheNextShort()
	}

	f.Close()
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func setupWebDriver() webdriver.Session {
	chromedriver := webdriver.NewChromeDriver("/Users/dooboolab/GoProjects/video-crawler/res/chromedriver")
	driverRunErr := chromedriver.Start()
	checkErr(driverRunErr)

	desired := webdriver.Capabilities{}
	required := webdriver.Capabilities{}

	session, err := chromedriver.NewSession(desired, required)
	checkErr(err)

	openURL := session.Url("https://www.youtube.com/")
	checkErr(openURL)
	session.SetTimeoutsImplicitWait(3000)

	return *session
}

func goToTheShort() {
	btn, err := session.FindElement(selenium.ByCSSSelector, "a[title='Shorts']")
	checkErr(err)

	btn.Click()
}

func getCurrnetShortUrl() string {
	time.Sleep(1000 * time.Millisecond)

	currentURL, err := session.GetUrl()
	checkErr(err)

	return currentURL
}

func goToTheNextShort() {
	nextButton, err := session.FindElement(selenium.ByCSSSelector, "button[aria-label='Next video']")
	checkErr(err)

	nextButton.Click()
}

func createTheCSV() (*os.File, *csv.Writer) {
	currentTime := time.Now()
	fileName := currentTime.Format("2006-01-02") + ".csv"
	f, err := os.Create(fileName)
	checkErr(err)

	w := csv.NewWriter(f)

	title := []string{"id", "url"}
	w.Write(title)
	w.Flush()

	return f, w
}
