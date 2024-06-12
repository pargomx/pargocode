package gecko

// Handler para robots.txt que permite todo.
func RobotsAllow(c *Context) error {
	return c.String(200, "User-agent: *\nAllow: /")
}

// Handler para robots.txt que restringe todo.
func RobotsDisallow(c *Context) error {
	return c.String(200, "User-agent: *\nDisallow: /")
}
