package main

import (

  "net/http"
  "github.com/go-redis/redis"
  "github.com/gin-gonic/gin"
  "log"
)

type Contact struct {
  Name string
  Number string
}


var router *gin.Engine 

var client *redis.Client

func main() {

  client = redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
  })

  router := gin.Default()

  router.Static("/css", "templates/css")
  router.LoadHTMLGlob("templates/*.html")
  router.GET("/", showIndexPage)
  router.GET("/add", showAddPage)
  router.POST("/save", saveContact)

  router.Run()

}

func showIndexPage(c *gin.Context) {

  contacts, err := client.LRange("contact",0,5).Result()

  if err != nil {
    log.Println("failed to connect database: contacts")
    log.Println("Error: ", err)
  }

  numbers, err := client.LRange("numeros",0,5).Result()

  if err != nil {
    panic("failed to connect database: numbers")
    log.Println("Error: ", err)
  }  

  c.HTML(
    http.StatusOK,
    "index.html",
    gin.H{
      "title":   "Agenda",
      "contacts": contacts,
      "numbers": numbers,
    },
  )
}

func showAddPage(c *gin.Context) {

  c.HTML(
    http.StatusOK,
    "add.html",
    gin.H{
      "title":   "Agregar",
    },
  )
 
}
 
func saveContact(c *gin.Context) {
// Save contact in DB
  var ct map[string]interface{}

  err := c.Bind(&ct)
  if err != nil {
    log.Printf("Couldn't decode data %+v\n", err)
  }

  log.Println("contact took")
  log.Printf("%+v\n", ct)

  name := c.PostFormArray("name")
  number:= c.PostFormArray("number")

  contact := client.LPush("contact", name)
  if contact != nil {
    log.Printf("LPush: %+v\n", contact)
  }

  numbers := client.LPush("numeros", number)
  if numbers != nil {
    log.Printf("LPush: %+v\n", numbers)
  }  

  log.Print("Contact Saved")

  c.HTML(
    http.StatusOK,
    "exito.html",
    gin.H{
      "title":   "Agenda",
    },
  )
}


