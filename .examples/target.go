package main

type Foo struct {
	SnakeAttribute          string `snake:"snakeAttribute"`
	CamelAttribute          string `camel:"camel-attribute"`
	PascalAttribute         string `pascal:"pascal_attribute"`
	KebabAttribute          string `kebab:"KebabAttribute"`
	ScreamingSnakeAttribute string `screamingSnake:"screamingSnake-Attribute"`
	ScreamingKebabAttribute string `screamingKebab:"screaming_kebabAttribute"`
	IgnoreAttribute         string `tagcase:"ignore" snake:"ignoreAttribute"`
	InitialismAttribute     string `camel:"blah_blah_blah_id"`
	MultiTagAttribute       string `snake:"multi_tag_attribute" camel:"MultiTagAttribute" pascal:"MultiTagAttribute"`
	WithDelimitAttribute    string `withDelimit:"omitzero,withDelimitAttribute,omitempty"`
}
