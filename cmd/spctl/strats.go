package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var spoLoginEndpoints = []string{
	"login.microsoftonline.com",
	"login.microsoftonline.de",
	"login.chinacloudapi.cn",
	"login-us.microsoftonline.com",
	"login-us.microsoftonline.com",
}

func resolveEnv(siteURL string) (string, error) {
	u, err := url.Parse(siteURL)
	if err != nil {
		return "", err
	}

	// Obvious SharePoint Online checks
	if strings.HasSuffix(u.Host, ".sharepoint.com") ||
		strings.HasSuffix(u.Host, ".sharepoint.de") ||
		strings.HasSuffix(u.Host, ".sharepoint.cn") ||
		strings.HasSuffix(u.Host, ".sharepoint-mil.us") ||
		strings.HasSuffix(u.Host, ".sharepoint.us") {
		return "spo", nil
	}

	// Check login redirect URL
	redirectURL, err := getRedirect(siteURL)
	if err != nil {
		return "", err
	}

	// fmt.Printf("Redirect URL: %s", redirectURL)

	u, err = url.Parse(redirectURL)
	if err != nil {
		return "", err
	}

	for _, endpoint := range spoLoginEndpoints {
		if strings.Contains(u.Host, endpoint) {
			return "spo", nil
		}
	}

	return "onprem", nil
}

func getStrategies(siteURL string) ([]string, error) {
	env, err := resolveEnv(siteURL)
	if err != nil {
		return nil, err
	}

	if env == "spo" {
		return []string{"ondemand", "azurecert", "azurecreds", "device", "saml", "addin"}, nil
	}

	redirectURL, err := getRedirect(siteURL)
	if err != nil {
		return nil, err
	}

	if isADFS(redirectURL) {
		return []string{"adfs", "ondemand"}, nil
	}

	if isFBA(redirectURL) {
		return []string{"fba", "ondemand"}, nil
	}

	if isTMG(redirectURL) {
		return []string{"tmg"}, nil
	}

	if isNTLM(redirectURL) {
		return []string{"ntlm", "ntlm2"}, nil
	}

	return []string{}, fmt.Errorf("can't resolve auth strategy")
}

func getResp(siteURL string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", siteURL, nil)
	if err != nil {
		return nil, err
	}

	// Pretend to be a browser
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	return resp, nil
}

func getRedirect(siteURL string) (string, error) {
	resp, err := getResp(siteURL)
	if err != nil {
		return "", err
	}

	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if redirect, err := resp.Location(); err == nil && redirect != nil {
		return getRedirect(redirect.String())
	}

	// fmt.Println(resp)
	return siteURL, nil
}

func isADFS(redirectURL string) bool {
	// https://fs.contoso.com/adfs/ls?version=1.0&action=signin&realm=urn%3AAppProxy%3Acom&appRealm=49c75c92-e1ee-eb11-94aa-02bfc0a8000c&returnUrl=https%3A%2F%2Fcontoso.com%2Flims&client-request-id=81D87280-3A04-0001-B426-AA8B043AD901%
	return strings.Contains(redirectURL, "/adfs/")
}

func isNTLM(siteURL string) bool {
	resp, err := getResp(siteURL)
	if err != nil {
		return false
	}

	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	/*
		&{401 Unauthorized 401 HTTP/1.1 1 1 map[Content-Length:[16] Content-Type:[text/plain; charset=utf-8] Date:[Sat, 15 Apr 2023 22:11:52 GMT] Microsoftsharepointteamservices:[16.0.0.4822] Request-Id:[6692a9a0-df1b-e01d-9b62-c7422321cf92] Server:[Microsoft-IIS/8.5 Microsoft-HTTPAPI/2.0] Spiislatency:[0] Sprequestduration:[2] Sprequestguid:[6692a9a0-df1b-e01d-9b62-c7422321cf92] Www-Authenticate:[NTLM Negotiate] X-Content-Type-Options:[nosniff] X-Frame-Options:[SAMEORIGIN] X-Ms-Invokeapp:[1; RequireReadOnly] X-Powered-By:[ASP.NET]] 0x14000160300 16 [] false false map[] 0x140000be600 0x140000ac160}
	*/

	// NTLM when WWW-Authenticate header is present and status code is 401
	return resp.StatusCode == http.StatusUnauthorized && resp.Header.Get("WWW-Authenticate") != ""
}

func isTMG(redirectURL string) bool {
	// https://contoso.com/CookieAuth.dll?GetLogon?curl=Z2F&reason=0&formdir=1%
	return strings.Contains(redirectURL, "/CookieAuth.dll")
}

func isFBA(redirectURL string) bool {
	// http://contoso.com/_layouts/15/spf/auth.aspx?ReturnUrl=%2fd%2fhr%2f_layouts%2f15%2fAuthenticate.aspx%3fSource%3d%252Fd%252Fhr&Source=%2Fd%2Fhr%
	return strings.Contains(redirectURL, "ReturnUrl=") && strings.Contains(redirectURL, "Authenticate.asp")
}
