package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
)

func getEnvVariable(key string) string {
	// Load the .env file in the current directory
	err := godotenv.Load(".env")
	if err !=nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx, cancel = chromedp.NewContext(ctx,
		// 设置日志方法
		chromedp.WithLogf(log.Printf),)

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)

	// 希望浏览器关闭，使用cancel()方法即可
	defer cancel()

	username := getEnvVariable("USERNAME")
	password := getEnvVariable("PASSWD")
	loginUrl := getEnvVariable("LOGIN_URL")

	err := chromedp.Run(ctx,
		login(loginUrl, username, password),
		ipoPage(),
	)
	if err != nil {
		log.Fatal(err)
	}	
}

// Login in tasks
func login(urlstr string, username string, password string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://www.rakuten-sec.co.jp`),
		// wait for login form element is visible (ie, page is loaded)
		chromedp.WaitVisible(`#form-login-id`, chromedp.ByID),
		// fill in username and password
		chromedp.SetValue(`#form-login-id`, username, chromedp.ByID),
		chromedp.SetValue(`#form-login-pass`, password, chromedp.ByID),
		// click login button
		chromedp.Click(`button[class="s3-form-login__btn"]`, chromedp.ByQuery),
	}
}

func ipoPage() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.WaitVisible(`#gmenu_domestic_stock`, chromedp.ByID),
		chromedp.Click(`#gmenu_domestic_stock`, chromedp.ByID),
		chromedp.Click(`ul[class="pcm-nav-03__list"] > li:nth-child(3)`, chromedp.ByQuery),
		chromedp.Sleep(20 * time.Second),
	}
}