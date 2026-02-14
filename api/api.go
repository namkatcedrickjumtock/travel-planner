package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)
 

func NewAPIListener(services any) (*gin.Engine, error) {
	route := gin.Default()
	route.Use(cors.Default())
 

	// route.POST("/users", func(ctx *gin.Context) {
	// 	var user *userModel.User
	// 	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{
	// 			"error": err.Error(),
	// 		})

	// 		return
	// 	}
	// 	user, err := usr.CreateUser(ctx, *user)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
	// 			Error: err.Error(),
	// 		})
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusOK, user)
	// })

	// route.GET("/users", func(ctx *gin.Context) {
	// 	userStartKey, err := strconv.ParseUint(ctx.Query("start_key"), 10, 32)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
	// 			Error: "Invalid start_key :" + err.Error(),
	// 		})
	// 		return
	// 	}
	// 	userCount, err := strconv.ParseUint(ctx.Query("count"), 10, 32)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
	// 			Error: `count's invalid` + err.Error(),
	// 		})
	// 		return
	// 	}
	// 	if userID := ctx.Query("user_id"); userID == "" {
	// 		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
	// 			Error: "user_id's required",
	// 		})
	// 		return
	// 	}

	// 	if userCount == 0 {
	// 		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
	// 			Error: "count must be greater than 1",
	// 		})
	// 		return
	// 	}

	// 	users, err := usr.GetAllUsers(ctx, int(userStartKey), int(userCount))
	// 	if err != nil {
	// 		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
	// 			Error: err.Error(),
	// 		})
	// 		return
	// 	}

	// 	ctx.JSON(http.StatusOK, users)
	// })
 
	return route, nil
}
