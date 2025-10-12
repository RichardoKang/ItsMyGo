package main

import (
	"gee"
	"net/http"
)

func main() {

	r := gee.New()

	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/helloGee", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello %s, you are at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	err := r.Run(":8080")
	if err != nil {
		return
	}

}
