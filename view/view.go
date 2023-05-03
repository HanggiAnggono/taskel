package view

import (
	"github.com/CloudyKit/jet"
	"github.com/gin-gonic/gin"
)

var JetView *jet.Set

func Init(r *gin.Engine) {
	JetView = jet.NewHTMLSet("./templates", "./templates/auth", "./templates/tasks", "./templates/globals")
}

func HTML(c *gin.Context, code int, template string, data map[string]interface{}, varmap jet.VarMap) {
	t, err := JetView.GetTemplate(template)
	if err != nil {
		panic(err)
	}

	c.Writer.WriteHeader(code)
	if err = t.Execute(c.Writer, varmap, data); err != nil {
		panic(err)
	}
}
