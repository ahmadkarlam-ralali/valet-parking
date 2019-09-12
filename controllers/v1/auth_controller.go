package v1

import (
	b64 "encoding/base64"

	"net/http"

	"github.com/ahmadkarlam-ralali/valet-parking/helpers"
	"github.com/ahmadkarlam-ralali/valet-parking/models"
	"github.com/ahmadkarlam-ralali/valet-parking/requests"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type AuthController struct {
	Db *gorm.DB
}

// Login godoc
// @Summary Login 
// @Description Login : Username : admin; Password : admin
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body requests.AuthLoginRequest true "Request Body"
// @Success 200 {string} string "Ok"
// @Failure 404 {string} string "Unauthorized"
// @Router /auth/login [post]
func (this *AuthController) Login(c *gin.Context) {
	var request requests.AuthLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.HttpError(c, "Wrong Format", http.StatusBadRequest)
		return
	}

	// check existancy
	var user models.User
	if err := this.Db.Where("username = ? and password = ?", request.Username, request.Password).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := b64.StdEncoding.EncodeToString([]byte(user.Username))

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
	})

}
