package controllers

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/naqeeb8a/Coffee-shop/initializers"
	"github.com/naqeeb8a/Coffee-shop/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var body struct {
		ProfilePicture  string
		FirstName       string
		LastName        string
		Mobile          string
		Email           string
		Password        string
		AppVersion      string
		DeviceOsVersion string
		DeviceModel     string
		DeviceUTCOffSet string
		DeviceToken     string
	}

	if (c.Bind(&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	var validationError = make(map[string]string)

	if body.ProfilePicture == "" {
		validationError["profilePicture"] = "Please provide a valid image"
	}
	if body.FirstName == "" {
		validationError["firstName"] = "Please provide a first name"
	}
	if body.LastName == "" {
		validationError["lastName"] = "Please provide a last name"
	}
	if body.Email == "" || !strings.Contains(body.Email, "@") || !strings.Contains(body.Email, ".") {
		validationError["email"] = "Please provide a valid email"
	}
	if len(body.Password) < 6 {
		validationError["password"] = "Password cannot be less than 6 characters"
	}
	if len(body.Mobile) < 11 {
		validationError["mobile"] = "Please provide a valid mobile number"
	}
	if len(validationError) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Fill all fields",
			"required": validationError,
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to Hash password",
		})
		return
	}

	user := models.User{Email: body.Email, Password: string(hash), Mobile: body.Mobile, FirstName: body.FirstName, LastName: body.LastName, FullName: body.FirstName + " " + body.LastName}
	result := initializers.DB.Create(&user)

	if (result.Error) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This user may exist already",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profilePicture": user.ProfilePicture,
		"firstName":      user.FirstName,
		"lastName":       user.LastName,
		"email":          user.Email,
		"mobile":         user.Mobile,
	})

}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if (c.Bind(&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	var validationError = make(map[string]string)
	if body.Email == "" || !strings.Contains(body.Email, "@") || !strings.Contains(body.Email, ".") {
		validationError["email"] = "Please provide a valid email"
		validationError["providedEmail"] = body.Email
	}
	if len(body.Password) < 6 {
		validationError["password"] = "Password cannot be less than 6 characters"
		validationError["providedPassword"] = body.Password
	}
	if len(validationError) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Fill all fields",
			"required": validationError,
		})
		return
	}
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No registered email found !!",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect password !!",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       user.ID,
		"expiresIn": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create Token !!",
		})
		return
	}
	user.AccessToken = tokenString
	user.ExpiryAt = time.Now().Add(time.Hour * 24 * 30).Unix()
	initializers.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"token":       tokenString,
		"expiresIn":   time.Now().Add(time.Hour * 24 * 30).Unix(),
		"userDetails": user,
	})

}

func EditUser(c *gin.Context) {
	var body struct {
		ProfilePicture string
		FirstName      string
		LastName       string
		Mobile         string
		Email          string
		Password       string
	}
	data, usererr := c.Get("user")
	if !usererr {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid token",
			"details": "No user assigned to this token please login again and pass the new token",
		})
		return
	}
	var loggedInUser models.User = data.(models.User)
	if (c.Bind(&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var user models.User
	initializers.DB.First(&user, loggedInUser.ID)
	if user.ID == 0 {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "No record found of user with this id",
		})
		return
	}
	if body.ProfilePicture != "" {
		user.ProfilePicture = body.ProfilePicture
	}
	if body.FirstName != "" {
		user.FirstName = body.FirstName
	}
	if body.LastName != "" {
		user.LastName = body.LastName
	}
	if body.Email != "" && !strings.Contains(body.Email, "@") && !strings.Contains(body.Email, ".") {
		user.Email = body.Email
	}

	if len(body.Mobile) >= 11 {
		user.Mobile = body.Mobile
	}
	initializers.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"profilePicture": user.ProfilePicture,
		"firstName":      user.FirstName,
		"lastName":       user.LastName,
		"email":          user.Email,
		"mobile":         user.Mobile,
	})

}

func GetUser(c *gin.Context) {
	data, _ := c.Get("user")
	var user models.User = data.(models.User)
	c.JSON(http.StatusOK, gin.H{"userDetails": user})

}
func AddUserLoyaltyPoints(c *gin.Context) {
	if c.Query("points") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":           "Query parameter (points) was not passed",
			"passed (points)": c.Query("points"),
		})
		return
	}

	points, err := strconv.Atoi(c.Query("points"))
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error":           "points should be integer",
			"passed (points)": c.Query("points"),
		})
		return
	}
	data, usererr := c.Get("user")
	if !usererr {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid token",
			"details": "No user assigned to this token please login again and pass the new token",
		})
		return
	}
	var loggedInUser models.User = data.(models.User)
	var user models.User
	initializers.DB.First(&user, loggedInUser.ID)
	if user.ID == 0 {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "No record found of user with this id",
		})
		return
	}
	user.LoyaltyPoints = points
	initializers.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"success": "loyalty points added successfully",
		"details": user,
	})
}
