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
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// // create headless chrome instance
	// ctx, cancel := chromedp.NewContext(context.Background())
	// defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	// 希望浏览器关闭，使用cancel()方法即可
	defer cancel()

	// login
	err :=login(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// get ipo participation links nodes
	joinNodes, applyNodes, err := getIpoNodes(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// do IPO
	err = doJoin(ctx, joinNodes)
	if err != nil {
		log.Fatal(err)
	}
	err = doApply(ctx, applyNodes)
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
		return fmt.Errorf("could not login: %s", err)
	}
	return nil
}

func getIpoNodes(ctx context.Context) ([]*cdp.Node, []*cdp.Node, error) {
	var joinNodes []*cdp.Node
	var applyNodes []*cdp.Node
	var ipoUrl string
	// navigate to ipo page
	if err := chromedp.Run(
		ctx,
		chromedp.WaitVisible(`#gmenu_domestic_stock`, chromedp.ByID),
		chromedp.Click(`#gmenu_domestic_stock`, chromedp.ByID),
		chromedp.Click(`ul.pcm-nav-03__list > li:nth-child(3)`, chromedp.ByQuery),
		chromedp.Nodes(`//div/nobr/a[contains(text(), "参加")]`, &joinNodes, chromedp.BySearch),
		chromedp.Nodes(`//div/a[contains(text(), "申込")]`, &applyNodes, chromedp.BySearch),
		chromedp.Location(&ipoUrl),
	); err != nil {
		return nil, nil, fmt.Errorf("could not get IPO nodes: %s", err)
	}
	return joinNodes,applyNodes, nil;
}

func doJoin(ctx context.Context, nodes []*cdp.Node) error {
	// force max timeout of 15 seconds for retrieving and processing the data
	shortPassword := getEnvVariable("SHORT_PASSWD")
	// iterate nodes
	for _, node := range nodes {
		if err := chromedp.Run(
			ctx,
			chromedp.MouseClickNode(node),
			chromedp.WaitVisible(`//input[@type='button' and contains(@value, "同意する") ]`, chromedp.BySearch),
			chromedp.Click(`//input[@type='button' and contains(@value, "同意する") ]`, chromedp.BySearch),
			// fill in ipo num (defalut value is 100)
			chromedp.WaitVisible(`//input[@type='submit' and contains(@value, "確　認") ]`, chromedp.BySearch),
			chromedp.SetValue(`//input[@type='text' and @name='value']`, "100", chromedp.BySearch),
			chromedp.SetValue(`#price`, "0", chromedp.ByID),
			chromedp.Click(`//input[@type='submit' and contains(@value, "確　認")]`, chromedp.BySearch),
			chromedp.Sleep(2*time.Second),
			// file in confirm passwd
			chromedp.WaitVisible(`//input[@type='submit' and contains(@value, "参加申込")]`, chromedp.BySearch),
			chromedp.SetValue(`//input[@type='password' and @name='password']`, shortPassword, chromedp.BySearch),
			chromedp.Click(`//input[@type='submit' and contains(@value, "参加申込")]`, chromedp.BySearch),
			// navigate back 3 times
			chromedp.WaitVisible(`//font[@class='ord_subtitle-01' and contains(text(), "受付完了") ]`, chromedp.BySearch),
			chromedp.NavigateBack(),
			chromedp.NavigateBack(),
			chromedp.NavigateBack(),
			chromedp.NavigateBack(),
		); err != nil {
			return fmt.Errorf("could not do ipo")
		}
	}
	return nil
}


func doApply(ctx context.Context, nodes []*cdp.Node) error {
	return nil
}