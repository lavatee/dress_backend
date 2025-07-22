package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lavatee/dresscode_backend/internal/model"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type SignUpInput struct {
	Name     string `json:"name" binding:"required" validate:"required,min=2,max=255"`
	Email    string `json:"email" binding:"required" validate:"required,email,max=255"`
	Password string `json:"password" binding:"required" validate:"required,min=8,max=72"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required" validate:"required,email,max=255"`
	Password string `json:"password" binding:"required" validate:"required"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (e *Endpoint) SignUp(c *gin.Context) {
	var input SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	if err := validate.Struct(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	userId, err := e.services.Auth.SignUp(model.User{Name: input.Name, Email: input.Email, Password: input.Password})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"id": userId})
}

func (e *Endpoint) SignIn(c *gin.Context) {
	var input SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	if err := validate.Struct(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	access, refresh, err := e.services.Auth.SignIn(input.Email, input.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"access_token": access, "refresh_token": refresh})
}

func (e *Endpoint) Refresh(c *gin.Context) {
	var input RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	access, refresh, err := e.services.Auth.Refresh(input.RefreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"access_token": access, "refresh_token": refresh})
}

type SetRoleInput struct {
	UserId int `json:"user_id" binding:"required"`
}

func (e *Endpoint) NewAdmin(c *gin.Context) {
	var input SetRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	if err = e.services.Auth.NewAdmin(userId, input.UserId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"message": "Admin set successfully"})
}

func (e *Endpoint) NewBuyer(c *gin.Context) {
	var input SetRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	if err = e.services.Auth.NewBuyer(userId, input.UserId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"message": "Buyer set successfully"})
}

func (e *Endpoint) RemoveBuyer(c *gin.Context) {
	var input SetRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	if err = e.services.Auth.RemoveBuyer(userId, input.UserId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"message": "Buyer removed successfully"})
}

func (e *Endpoint) GetUserRole(c *gin.Context) {
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	role, err := e.services.Auth.GetUserRole(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"role": role})
}

func (e *Endpoint) GetUser(c *gin.Context) {
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	user, err := e.services.Auth.GetUser(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"user": user})
}

func (e *Endpoint) GetUserIdByToken(c *gin.Context) {
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"user_id": userId})
}
