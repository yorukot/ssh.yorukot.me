package content

type Content struct {
	HeaderTitle   string
	HeaderTagline string
	FooterQuote   string
	FooterLinks   []FooterLink
}

type FooterLink struct {
	Label string
	URL   string
}

func GetContent() Content {
	return Content{
		HeaderTitle:   "Yorukot",
		HeaderTagline: "Open-source developer",
		FooterQuote:   "Get busy living, or get busy dying",
		FooterLinks: []FooterLink{
			{Label: "Email", URL: "mailto:hi@yorukot.me"},
			{Label: "GitHub", URL: "https://github.com/yorukot"},
			{Label: "Telegram", URL: "https://t.me/yorukot"},
			{Label: "Discord", URL: "https://dc.yorukot.me"},
			{Label: "Ko-fi", URL: "https://donate.yorukot.me"},
		},
	}
}

func HomePage() string {
	
	return ""
}