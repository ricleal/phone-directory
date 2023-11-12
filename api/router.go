package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"phone-directory/internal/entities"
	"phone-directory/internal/service"
)

func Router(h *handler) (*gin.Engine, error) {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	v1 := router.Group("/v1")
	v1.GET("/users/:id", h.getUser)
	v1.POST("/users", h.postUser)
	v1.PUT("/users/:id", h.putUser)
	v1.DELETE("/users/:id", h.deleteUser)

	return router, nil
}

type handler struct {
	us service.UserService
	ps service.PhoneService
	as service.AddressService
}

func NewHandler(us service.UserService, ps service.PhoneService, as service.AddressService) *handler {
	return &handler{
		us: us,
		ps: ps,
		as: as,
	}
}

func (h *handler) getUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.us.Get(c.Request.Context(), uint(idInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *handler) postUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newUser := entities.User{
		Name: user.Name,
	}
	for _, phone := range user.Phones {
		newUser.Phones = append(newUser.Phones, entities.Phone{
			Number: phone,
		})
	}
	for _, address := range user.Addresses {
		newUser.Addresses = append(newUser.Addresses, entities.Address{
			Address: address,
		})
	}
	h.us.Create(c.Request.Context(), &newUser)
	c.JSON(http.StatusCreated, newUser)
}

func (h *handler) putUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newUser := entities.User{
		ID:   uint(idInt),
		Name: user.Name,
	}
	for _, phone := range user.Phones {
		newUser.Phones = append(newUser.Phones, entities.Phone{
			Number: phone,
		})
	}
	for _, address := range user.Addresses {
		newUser.Addresses = append(newUser.Addresses, entities.Address{
			Address: address,
		})
	}
	h.us.Update(c.Request.Context(), &newUser)
	c.JSON(http.StatusOK, newUser)
}

func (h *handler) deleteUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.us.Delete(c.Request.Context(), uint(idInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}
