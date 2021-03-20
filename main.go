package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
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
	// opts := append(chromedp.DefaultExecAllocatorOptions[:],
	// 	chromedp.Flag("headless", false),
	// )
	// ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	// defer cancel()

	// ctx, cancel = chromedp.NewContext(ctx)
	// defer cancel()

	// create headless chrome instance
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	// 希望浏览器关闭，使用cancel()方法即可
	defer cancel()

	// login
	err :=login(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// get ipo links nodes
	joinNodes, err := ipoPage(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, node := range joinNodes {
		fmt.Println(node.AttributeValue("href"))
	}
}

// Login in tasks
func login(ctx context.Context) error {
	username := getEnvVariable("USERNAME")
	password := getEnvVariable("PASSWD")
	loginUrl := getEnvVariable("LOGIN_URL")
	// login 
	if err := chromedp.Run(
		ctx,
		chromedp.Navigate(loginUrl),
		// wait for login form element is visible (ie, page is loaded)
		chromedp.WaitVisible(`#form-login-id`, chromedp.ByID),
		// fill in username and password
		chromedp.SetValue(`#form-login-id`, username, chromedp.ByID),
		chromedp.SetValue(`#form-login-pass`, password, chromedp.ByID),
		// click login button
		chromedp.Click(`button.s3-form-login__btn`, chromedp.ByQuery),
	); err != nil {
		return fmt.Errorf("could not login %s", err)
	}
	return nil
}

func ipoPage(ctx context.Context, ) ([]*cdp.Node, error) {
	var nodes []*cdp.Node
	// navigate to ipo page
	if err := chromedp.Run(
		ctx,
		chromedp.WaitVisible(`#gmenu_domestic_stock`, chromedp.ByID),
		chromedp.Click(`#gmenu_domestic_stock`, chromedp.ByID),
		chromedp.Click(`ul.pcm-nav-03__list > li:nth-child(3)`, chromedp.ByQuery),
		chromedp.Nodes(`//div/nobr/a[contains(text(), "参加")]`, &nodes, chromedp.BySearch),
	); err != nil {
		return nil, fmt.Errorf("could not login %s", err)
	}
	return nodes, nil;
}