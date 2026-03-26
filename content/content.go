package content

type Content struct {
	HeaderTitle string
}

func GetContent() Content {
	return Content{
		HeaderTitle: "Yorukot",
	}
}