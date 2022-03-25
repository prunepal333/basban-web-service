package main
import (
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
	"encoding/json"
	"io/ioutil"
)

type Book struct {
	ISBN string	`json:"isbn"`
	Title string `json:"title"`
	Price float64 `json:"price"`
	OwnerId int `json:"owner_id"`
	Description string `json:"description"`
	ImageURI string `json:"image_uri"`
}

type Owner struct {
	Id int `json:"id"`
	Name string `json:"name"`
	phone string `json:"phone"`
	email string `json:"email"`
}

var books []Book
var owners []Owner

func main() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	booksJson, _ := ioutil.ReadFile("books.json")
	ownersJson, _ := ioutil.ReadFile("owners.json")

	json.Unmarshal([]byte(booksJson), &books)
	json.Unmarshal([]byte(ownersJson), &owners)


	router.Use(CORSMiddleware())
	router.GET("/", error404)
	router.GET("/books", getBooks)
	router.GET("/books/:slug", getBookById)
	router.GET("/owners", getOwners)
	router.GET("/owners/:slug", getOwnerById)
	router.GET("/owners/:slug/books", getBooksByOwner)

	router.Run("localhost:5000")
}
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}
func getBookById(c *gin.Context) {
	id := c.Param("slug");
	for _, a := range books {
		if a.ISBN == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not Found"})
}
func error404(c *gin.Context) {
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not Found"})
}
func getOwners(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, owners)
}
func getOwnerById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("slug"));
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not Found"})
	}
	for _, o := range owners {
		if o.Id == id {
			c.IndentedJSON(http.StatusOK, o)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not Found"})
}
func getBooksByOwner(c *gin.Context) {
	owner_id, _ := strconv.Atoi(c.Param("slug"))
	if(!isValidOwner(owner_id)){
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not Found"})
	}
	booksByOwner := []Book{}
	for _, b := range books {
		if b.OwnerId == owner_id {
			booksByOwner = append(booksByOwner, b)
		}
	}
	c.IndentedJSON(http.StatusOK, booksByOwner)
}
func isValidOwner(id int) bool{
	for _, o := range owners {
		if o.Id == id{
			return true
		}
	}
	return false
}
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {

        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "*")
        c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}