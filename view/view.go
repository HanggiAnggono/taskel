package view

import (
	"github.com/CloudyKit/jet"
	"github.com/gin-gonic/gin"
)

var JetView *jet.Set

func Init(r *gin.Engine) {
	JetView = jet.NewHTMLSet("./templates", "./templates/auth", "./templates/globals")
}

func HTML(c *gin.Context, code int, template string, data map[string]interface{}) {
	t, err := JetView.GetTemplate(template)
	if err != nil {
		panic(err)
	}

	vars := make(jet.VarMap)
	for key, v := range data {
		vars.Set(key, v)
	}
	c.Writer.WriteHeader(code)
	if err = t.Execute(c.Writer, vars, vars); err != nil {
		panic(err)
	}
}
