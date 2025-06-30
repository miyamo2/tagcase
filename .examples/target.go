package main

type Foo struct {
	SnakeAttribute          string `snake:"snake_attribute"`
	CamelAttribute          string `camel:"camelAttribute"`
	PascalAttribute         string `pascal:"PascalAttribute"`
	KebabAttribute          string `kebab:"kebab-attribute"`
	ScreamingSnakeAttribute string `screamingSnake:"SCREAMING_SNAKE_ATTRIBUTE"`
	ScreamingKebabAttribute string `screamingKebab:"SCREAMING-KEBAB-ATTRIBUTE"`
	IgnoreAttribute         string `tagcase:"ignore" snake:"ignoreAttribute"`
	InitialismAttribute     string `camel:"somethingID"`
	MultiTagAttribute       string `snake:"multi_tag_attribute" camel:"multiTagAttribute" pascal:"MultiTagAttribute"`
}
