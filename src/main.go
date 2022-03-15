package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fedesog/webdriver"
	"github.com/joho/godotenv"
	"github.com/tebeka/selenium"
)

var contries []string
var session webdriver.Session

func main() {
	contries = []string{"Japan", "United Kingdom", "South Korea", "United States"}
	videoCounts := 30
	godotenv.Load("/Users/dooboolab/GoProjects/video-crawler/.env")

	session = setupWebDriver()

	goToTheSignInPage()
	signIn()

	// fileName := createTheCSV()
	// f, err := os.Create(fileName)
	// checkErr(err)

	f, w := createTheCSV()

	for _, contry := range contries {
		selectTheContry(contry)

		time.Sleep(4 * time.Second)
		session.Refresh()
		time.Sleep(1 * time.Second)

		goToTheShort()

		for i := 1; i <= videoCounts; i++ {
			currentURL := getCurrnetShortUrl()
			fmt.Println(currentURL)
			data := []string{currentURL, currentURL, contry}

			w.Write(data)
			w.Flush()

			time.Sleep(1 * time.Second)
			goToTheNextShort()
		}
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

func goToTheSignInPage() {
	signInButton, err := session.FindElement(selenium.ByCSSSelector, "tp-yt-paper-button[aria-label='Sign in']")
	checkErr(err)
	signInButton.Click()
}

func signIn() {
	inputEmail, err := session.FindElement(selenium.ByCSSSelector, "input[type='email']")
	checkErr(err)

	email := os.Getenv("USER_ID")
	inputEmail.SendKeys(email)

	nextButton, err := session.FindElement(selenium.ByCSSSelector, "#identifierNext")
	checkErr(err)
	nextButton.Click()

	time.Sleep(2 * time.Second)
	inputPassword, err := session.FindElement(selenium.ByCSSSelector, "input[type='password']")
	checkErr(err)
	password := os.Getenv("USER_PW")
	inputPassword.SendKeys(password)

	signInButton, err := session.FindElement(selenium.ByCSSSelector, "#passwordNext")
	checkErr(err)
	signInButton.Click()
	time.Sleep(1 * time.Second)
}

func selectTheContry(country string) {
	profileButton, err := session.FindElement(selenium.ByCSSSelector, "#avatar-btn")
	checkErr(err)
	profileButton.Click()

	locationButton, err := session.FindElements(selenium.ByCSSSelector, "#endpoint > tp-yt-paper-item")
	checkErr(err)
	locationButton[6].Click()

	locations, err := session.FindElements(selenium.ByCSSSelector, "#items #endpoint yt-formatted-string[id='label']")
	checkErr(err)
	index, isFound := findTheCuntryIndex(locations, country)
	if isFound {
		locations[index].Click()
	}
}

func findTheCuntryIndex(locations []webdriver.WebElement, country string) (int, bool) {
	for i, location := range locations {
		text, err := location.Text()
		checkErr(err)
		fmt.Println(i, text)
		if text == country {
			return i, true
		}
	}

	return -1, false
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

	title := []string{"id", "url", "state"}
	w.Write(title)
	w.Flush()

	return f, w
}
