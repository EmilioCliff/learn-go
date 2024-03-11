package server

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var users []User

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
}

func (server *Server) login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, user := range users {
		if user.Username == req.Username {
			if user.Password == req.Password {
				accessToken, err := server.tokenMaker.CreateToken(req.Username, time.Minute)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				rsp := loginResponse{
					AccessToken: accessToken,
					User:        user,
				}
				ctx.JSON(http.StatusOK, rsp)
				return
			}
			ctx.JSON(http.StatusForbidden, gin.H{"error": "incorrect password"})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, users)
	return

}

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = strconv.Itoa(rand.Intn(1000))
	users = append(users, user)

	ctx.JSON(http.StatusOK, users)
	return
}

type deteleUserRequet struct {
	ID string `uri:"id" binding:"required"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req deteleUserRequet
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for idx, user := range users {
		if user.ID == req.ID {
			users = append(users[:idx], users[idx+1:]...)
			ctx.JSON(http.StatusOK, users)
			return
		}
	}
	ctx.JSON(http.StatusNotFound, users)
	return
}
