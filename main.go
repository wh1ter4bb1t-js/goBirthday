package main

import (
	"fmt"
	"log"

	"github.com/mxschmitt/playwright-go"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config/goBirthday/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	emailAdd := viper.GetString("email")
	pass := viper.GetString("password")

	pw, err := playwright.Run()
	checkErr(err, "failed to start to start playwright")

	browser, err := pw.Chromium.Launch()
	checkErr(err, "failed to launch chromium browser")

	page, err := browser.NewPage()
	checkErr(err, "failed to open new page in browser")

	if _, err = page.Goto("https://facebook.com"); err != nil {
		log.Fatal("failed to open login page")
	}

	email, err := page.QuerySelectorAll("#email")
	checkErr(err, "failed to get email input field")

	password, err := page.QuerySelectorAll("#pass")
	checkErr(err, "failed to get password input field")

	submit, err := page.QuerySelector("data-testid=royal_login_button")
	checkErr(err, "failed to get submit button")

	err = email[0].Fill(emailAdd)
	checkErr(err, "failed to fill in email")

	err = password[0].Fill(pass)
	checkErr(err, "failed to fill in password")

	submit.Click()
	page.WaitForTimeout(4000)
	if _, err = page.Goto("https://www.facebook.com/events/birthdays/"); err != nil {
		log.Fatal("failed to open birthdays page")
	}

	page.Click("main")

	recentInputs, err := page.QuerySelectorAll("._5rpu")
	checkErr(err, "failed to get birthday message inputs")

	for _, inputs := range recentInputs {
		inputs.Fill("Happy Birthday! I hope you have a good one!")
		err = page.Keyboard().Press("Enter")
		checkErr(err, "failed to press enter key")
	}
}

func checkErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:\nerror: %s", message, err)
	}
}
