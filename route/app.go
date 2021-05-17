package route

import (
	"github.com/gin-gonic/gin"

	"goblog/controller"

	// "net/http"
	// "strconv"
)


func Router(g *gin.Engine) {

	g.GET("/", controller.Index)

	g.GET("/page/:p", controller.Index)

	g.GET("/category/:name", controller.Category)
	g.GET("/category/:name/page/:p", controller.Category)

	g.GET("/archive/:y/:m", controller.Archive)
	g.GET("/archive/:y/:m/page/:p", controller.Archive)

	g.GET("/archives/:url", controller.Article)

	g.GET("about.html", controller.Page)

	

	// m.Get("/search/:k", func(r render.Render, params martini.Params) {
	// 	var controller = controller.NewSearchController()
	// 	controller.Search(r, params["k"])
	// })

	// m.Get("/search/:k/page/:p", func(r render.Render, params martini.Params) {
	// 	var controller = controller.NewSearchController()
	// 	var p int
	// 	p, _ = strconv.Atoi(params["p"])
	// 	controller.Search(r, params["k"], p)
	// })


	g.NoRoute(controller.NotFound)

}
