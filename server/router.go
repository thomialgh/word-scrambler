package server

import (
	"io"
	"net/http"
	"text/template"
	"word-scrambler/controllers"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func router() http.Handler {
	r := echo.New()
	r.Static("/static", "assets")
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	r.Renderer = renderer
	rAPI := r.Group("/api")
	{
		rAPI.GET("/question", controllers.GetQuestion, controllers.AuthHandler)
		rAPI.POST("/question/:question_id/answer", controllers.Answers, controllers.AuthHandler)
		rAPI.POST("/login", controllers.Login)
		rAPI.GET("/protected", controllers.Protected, controllers.AuthHandler)
		rAPI.GET("/total-score", controllers.GetTotalScore, controllers.AuthHandler)
		rAPI.GET("/summaries", controllers.Summary, controllers.AuthHandler)
	}

	r.GET("/login", controllers.LoginPage)
	r.GET("", controllers.Index)

	return r
}
