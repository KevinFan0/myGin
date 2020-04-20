// Custom Validators
package main

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"fmt"
)

// Booking contains binded and validated data.
type Booking struct {
	CheckIn		time.Time	`form:"check_in" binding:"required" time_format:"2006-01-02"`
	CheckOut	time.Time	`form:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

func getBookable(c *gin.Context) {
	var b Booking
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// Only Bind Query String
type Person struct {
	Name		string		`form:"name"`
	Address		string		`form:"address"`
	Birthday	time.Time	`form:"birthday" time_format:"2006-01-02" time_utc:"1"`
	CreateTime	time.Time	`form:"createTime" time_format:"unixNano"`
	UnixTime	time.Time	`form:"unixTime" time_format:"unix"`
}


func startPage(c *gin.Context)  {
	var person Person
	if c.ShouldBindQuery(&person) == nil {
		log.Println("====== Only Bind By Query String ======")
		log.Println(person.Name)
		log.Println(person.Address)
	}
	c.JSON(200, "success")
}

// Bind Query String or Post Data
func startPage2(c *gin.Context)  {
	var person Person
	// If `GET`, only `Form` binding engine (`query`) used.
	// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if c.ShouldBind(&person) == nil {
		log.Println(person.Name)
		log.Println(person.Address)
		log.Println(person.Birthday)
		log.Println(person.CreateTime)
		log.Println(person.UnixTime)
	}
	c.JSON(200, "success")
}


// Bind Uri
type PersonInfo struct {
	ID		string	`uri:"id" binding:"required,uuid"`
	Name	string	`uri:"name" binding:"required"`
}

// Bind Header
type testHeader struct {
	Rate	int		`header:"Rate"`
	Domain	string	`header:"Domain"`
}

// Bind HTML checkboxes
type myForm struct {
	Colors	[]string	`form:"colors[]"`
}

func formHandler(c *gin.Context)  {
	var fakeForm myForm
	c.ShouldBind(&fakeForm)
	c.JSON(200, gin.H{"color": fakeForm.Colors})
}


// Multipart/Urlencoded binding
type ProfileForm struct {
	Name	string					`form:"name" binding:"required"`
	Avatar	*multipart.FileHeader	`form:"avatar" binding:"required"`
	// or for multiple files
	// Avatars []*multipart.FileHeader `form:"avatar" binding:"required"`
}

func main()  {
	route := gin.Default()
	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("bookabledate", bookableDate)
	// }
	// route.GET("/bookable", getBookable)
	// route.Any("/testing", startPage)
	// route.GET("/testing2", startPage2)
	route.GET("/:name/:id", func(c *gin.Context)  {
		var person PersonInfo
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
	})

	route.GET("/", func(c *gin.Context) {
		h := testHeader{}
		if err := c.ShouldBind(&h); err != nil {
			c.JSON(200, err)
		}
		fmt.Printf("%#v\n", h)
		c.JSON(200, gin.H{"Rate":h.Rate, "Domain":h.Domain})
	})

	route.POST("/profile", func(c *gin.Context) {
		// you can bind multipart form with explicit binding declaration:
		// c.ShouldBindWith(&form, binding.Form)
		// or you can simply use autobinding with ShouldBind method:
		var form ProfileForm
		// in this case proper binding will be automatically selected
		if err := c.ShouldBind(&form); err != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}
		err := c.SaveUploadedFile(form.Avatar, form.Avatar.Filename)
		if err != nil {
			c.String(http.StatusInternalServerError, "unknown error")
			return
		}
		c.String(http.StatusOK, "OK")
	})

	// XML, JSON, YAML and ProtoBuf rendering
	// gin.H is a shortcut for map[string]interface{}
	route.GET("/someJSON", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	route.GET("/moreJSON", func(c *gin.Context){
		var msg struct {
			Name	string	`json:"user"`
			Message	string
			Number	int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		// Note that msg.Name becomes "user" in the JSON
		// Will output  :   {"user": "Lena", "Message": "hey", "Number": 123}
		c.JSON(http.StatusOK, msg)
	})

	route.GET("/someXML", func(c *gin.Context){
		c.XML(http.StatusOK, gin.H{"message": "hey", "status": })
	})

	route.Run(":8082")
}