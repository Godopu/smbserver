package router

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"smbserver/model"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	apiEngine := gin.New()
	apiv1 := apiEngine.Group("api/v1")
	{
		apiv1.POST("/alarm", PostAlarm)
		apiv1.GET("/alarm", GetAlarm)
	}
	assetEngine := gin.New()
	assetEngine.Static("/", "./front/build/web")

	r := gin.New()
	r.Any("/*any", func(c *gin.Context) {
		path := c.Param("any")
		if strings.HasPrefix(path, "/api/v1") {
			apiEngine.HandleContext(c)
		} else {
			assetEngine.HandleContext(c)
		}
	})

	return r
}

func handleError(c *gin.Context) {
	if r := recover(); r != nil {
		log.Println(r)
		c.String(http.StatusBadRequest, r.(error).Error())
	}
}

func PostAlarm(c *gin.Context) {
	// hms
	defer handleError(c)
	var aTime map[string]interface{}
	err := c.BindJSON(&aTime)
	if err != nil {
		panic(err)
	}

	n := time.Now()

	h := int(aTime["h"].(float64))
	m := int(aTime["m"].(float64))

	if h > 24 || h < 0 {
		panic(errors.New("wrong hour is entered"))
	}

	if m > 60 || m < 0 {
		panic(errors.New("wrong minute is entered"))
	}

	at := time.Date(n.Year(), n.Month(), n.Day(), h, m, 0, 0, n.Location())
	fmt.Println(at)
	if at.Sub(n).Seconds() <= 0 {
		panic(errors.New("please enter the time after now on."))
	}

	model.Alarm = append(model.Alarm, at)
	c.String(http.StatusCreated, "OK")
}

func GetAlarm(c *gin.Context) {
	// hms
	defer handleError(c)
	builder := strings.Builder{}
	for _, v := range model.Alarm {
		builder.Write([]byte(fmt.Sprintf("%2d:%2d\n", v.Hour(), v.Minute())))
	}

	c.String(http.StatusOK, builder.String())
}
