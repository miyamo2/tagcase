package main

type Foo struct {
	SnakeAttribute          string `snake:"snake_Attribute"`
	CamelAttribute          string `camel:"CamelAttribute"`
	PascalAttribute         string `pascal:"pascalAttribute"`
	KebabAttribute          string `kebab:"kebab_attribute"`
	ScreamingSnakeAttribute string `screamingSnake:"ScreamingSnakeAttribute"`
	ScreamingKebabAttribute string `screamingKebab:"ScreamingKebabAttribute"`
	IgnoreAttribute         string `tagcase:"ignore" snake:"ignoreAttribute"`
	InitialismAttribute     string `camel:"somethingID"`
	MultiTagAttribute       string `snake:"multi_tag_attribute" camel:"multiTagAttribute" pascal:"MultiTagAttribute"`
}
