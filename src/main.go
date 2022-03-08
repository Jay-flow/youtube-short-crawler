package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fedesog/webdriver"
	"github.com/tebeka/selenium"
)

func main() {
	videoNumber := 100
	session := setupWebDriver()
	goToTheShort(session)

	for i := 1; i <= videoNumber; i++ {
		currentURL := getCurrnetShortUrl(session)
		fmt.Println(currentURL)

		time.Sleep(1 * time.Second)
		goToTheNextShort(session)
	}
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

	return *session
}

func goToTheShort(session webdriver.Session) {
	openURL := session.Url("https://www.youtube.com/")
	checkErr(openURL)

	session.SetTimeoutsImplicitWait(3000)
	btn, err := session.FindElement(selenium.ByCSSSelector, "a[title='Shorts']")
	checkErr(err)

	btn.Click()
}

func getCurrnetShortUrl(session webdriver.Session) string {
	time.Sleep(1000 * time.Millisecond)

	currentURL, err := session.GetUrl()
	checkErr(err)

	return currentURL
}

func goToTheNextShort(session webdriver.Session) {
	nextButton, err := session.FindElement(selenium.ByCSSSelector, "button[aria-label='Next video']")
	checkErr(err)

	nextButton.Click()
}
