package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"xks/analyticrelationship/util"
	"xks/analyticrelationship/util/logger"
	"xks/analyticrelationship/util/web"
)

func main() {
	var url string

	flag.StringVar(&url, "u", "", "URL to extract Google Analytics ID")
	flag.StringVar(&url, "url", "", "URL to extract Google Analytics ID")

	const usage = `Usage: ./analyticrelationship -u [TARGET]`

	flag.Usage = func() {
		fmt.Println(util.PrintBanner())
		fmt.Println(usage)
	}
	flag.Parse()

	if url != "" {
		web.Start(url)
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				if err := scanner.Err(); err != nil {
					logger.Crash("bufio couldn't read stdin correctly.", err)
				} else {
					web.Start(scanner.Text())
				}
			}
		}
	}
}
