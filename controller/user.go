package controller

import (
	"net/http"

	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	//1. deal with the parameters
	p := new(models.ParamSignUp)

	//the request parameters Error (Formatting Error)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)),
		})
		return
	}

	//Deprecates: Estimate the param with business logic (Manually)
	//if len(p.Password) == 0 || len(p.RePassword) == 0 || len(p.Username) == 0 || p.Password != p.RePassword {
	//	zap.L().Error("SignUp with invalid params")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "logic of request parameters Error",
	//	})
	//	return
	//}

	//2. deal with the business logic
	logic.SignUp(p)

	//3. return the response
}
