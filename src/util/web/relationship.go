package web

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
	"xks/analyticrelationship/util/logger"
)

func GetGoogleTagManager(targetURL string) (bool, []string) {
	response, _ := getURLResponse(targetURL)
	var resultTagManager []string

	if response == "" {
		return false, resultTagManager
	}

	patterns := []struct {
		regex       *regexp.Regexp
		replaceFunc func(string) string
	}{
		{regexp.MustCompile(`www\.googletagmanager\.com/ns\.html\?id=[A-Z0-9\-]+`), func(s string) string {
			return "https://" + strings.Replace(s, "ns.html", "gtm.js", -1)
		}},
		{regexp.MustCompile("GTM-[A-Z0-9]+"), func(s string) string {
			return "https://www.googletagmanager.com/gtm.js?id=" + s
		}},
		{regexp.MustCompile(`UA-\d+-\d+`), nil},
	}

	for _, pattern := range patterns {
		data := pattern.regex.FindStringSubmatch(response)
		if len(data) > 0 {
			if pattern.replaceFunc != nil {
				resultTagManager = append(resultTagManager, pattern.replaceFunc(data[0]))
			} else {
				aux := pattern.regex.FindAllStringSubmatch(response, -1)
				var result []string
				for _, r := range aux {
					result = append(result, r[0])
				}
				return true, result
			}
		}
	}

	return false, resultTagManager
}

func getUA(url string) []string {
	pattern := regexp.MustCompile("UA-[0-9]+-[0-9]+")
	response, _ := getURLResponse(url)

	if response == "" {
		return nil
	}

	aux := pattern.FindAllStringSubmatch(response, -1)
	result := make([]string, len(aux))
	for i, r := range aux {
		result[i] = r[0]
	}

	return result
}

func getURLResponse(url string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 3,
	}
	res, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func getDomainsFromHackerTarget(id string) []string {
	url := fmt.Sprintf("https://api.hackertarget.com/analyticslookup/?q=%s", id)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("HTTP request failed with status:", resp.Status)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	lines := strings.Split(string(body), "\n")
	if len(lines) > 0 && !strings.Contains(lines[0], "API count exceeded") {
		return lines
	}

	return nil
}

func getDomains(id string) []string {
	allDomains := []string{}
	domains := getDomainsFromHackerTarget(id)
	uniqueDomains := make(map[string]bool)

	for _, domain := range domains {
		if !uniqueDomains[domain] {
			uniqueDomains[domain] = true
			allDomains = append(allDomains, domain)
		}
	}

	return allDomains
}

func showDomains(ua string) {
	allDomains := getDomains(ua)
	fmt.Println("~> " + ua)
	if len(allDomains) == 0 {
		logger.Warning("NOT FOUND")
	}
	for _, domain := range allDomains {
		if strings.Contains(domain, "error getting results") {
			printTree("No Relationships", 1)
		} else {
			printTree(domain, 1)
		}
	}
	fmt.Println("")
}

func contains(data []string, value string) bool {
	for _, v := range data {
		if v == value {
			return true
		}
	}
	return false
}

func Start(url string) {
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	uaResult, resultTagManager := GetGoogleTagManager(url)
	if len(resultTagManager) > 0 {
		var visited = []string{}
		var allUAs []string
		if !uaResult {
			urlGoogleTagManager := resultTagManager[0]
			logger.Info("Found ID: " + strings.Split(urlGoogleTagManager, "?id=")[1])
			allUAs = getUA(urlGoogleTagManager)
		} else {
			logger.Info("Found UA directly")
			allUAs = resultTagManager
		}
		logger.Info("Searching for relationships...\n")
		for _, ua := range allUAs {
			baseUA := strings.Join(strings.Split(ua, "-")[0:2], "-")
			if !contains(visited, baseUA) {
				visited = append(visited, baseUA)
				showDomains(baseUA)
			}
		}
		logger.Info("Finished!")
	} else {
		logger.Warning("Tagmanager URL not found")
	}
}

func printTree(node string, depth int) {
	prefix := strings.Repeat("  ", depth)
	fmt.Println(prefix + "\\__ " + node)
}
