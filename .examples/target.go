package main

type Foo struct {
	SnakeAttribute          string `snake:"snake_Attribute"`
	CamelAttribute          string `camel:"CamelAttribute"`
	PascalAttribute         string `pascal:"pascalAttribute"`
	KebabAttribute          string `kebab:"kebab_attribute"`
	ScreamingSnakeAttribute string `screamingSnake:"ScreamingSnakeAttribute"`
	ScreamingKebabAttribute string `screamingKebab:"ScreamingKebabAttribute"`
	IgnoreAttribute         string `tagcase:"ignore" snake:"ignoreAttribute"`
	InitialismAttribute     string `camel:"something_id"`
	MultiTagAttribute       string `snake:"multi-tag-attribute" camel:"MultiTagAttribute" pascal:"multiTagAttribute"`
}
