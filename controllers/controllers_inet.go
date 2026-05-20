package controllers

import (
	"log"
	"strconv"
	"strings"

	"go-fiber-test/database"
	m "go-fiber-test/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func HelloTest(c * fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func HelloTestV2(c * fiber.Ctx) error {
	return c.SendString("Hello, World! V2")
}

func BodyPersonTest(c *fiber.Ctx) error {
	p := new(m.Person)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	log.Println(p.Name) // john
	log.Println(p.Pass) // doe
	str := p.Name + p.Pass
	return c.JSON(str)
}

func ParamsTest(c *fiber.Ctx) error {

	str := "hello ==> " + c.Params("name")
	return c.JSON(str)
}

func QueryTest(c *fiber.Ctx) error {
	a := c.Query("search")
	str := "my search is  " + a
	return c.JSON(str)
}

func ValidTest(c *fiber.Ctx) error {
	//Connect to database

	user := new(m.User)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	validate := validator.New()
	errors := validate.Struct(user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.Error())
	}
	return c.JSON(user)
}	

func GetDogs(c *fiber.Ctx) error {
   db := database.DBConn
   var dogs []m.Dogs

   db.Find(&dogs) //delelete = null
   return c.Status(200).JSON(dogs)
}

func GetDog(c *fiber.Ctx) error {
   db := database.DBConn
   search := strings.TrimSpace(c.Query("search"))
   var dog []m.Dogs

   result := db.Find(&dog, "dog_id = ?", search)

   // returns found records count, equals `len(users)
   if result.RowsAffected == 0 {
       return c.SendStatus(404)
   }
   return c.Status(200).JSON(&dog)
}

func AddDog(c *fiber.Ctx) error {
   //twst3
   db := database.DBConn
   var dog m.Dogs

   if err := c.BodyParser(&dog); err != nil {
       return c.Status(503).SendString(err.Error())
   }

   db.Create(&dog)
   return c.Status(201).JSON(dog)
}

func UpdateDog(c *fiber.Ctx) error {
   db := database.DBConn
   var dog m.Dogs
   id := c.Params("id")

   if err := c.BodyParser(&dog); err != nil {
       return c.Status(503).SendString(err.Error())
   }

   db.Where("id = ?", id).Updates(&dog)
   return c.Status(200).JSON(dog)
}

func RemoveDog(c *fiber.Ctx) error {
   db := database.DBConn
   id := c.Params("id")
   var dog m.Dogs

   result := db.Delete(&dog, id)

   if result.RowsAffected == 0 {
       return c.SendStatus(404)
   }

   return c.SendStatus(200)
}

func GetDogsJson(c *fiber.Ctx) error {
   db := database.DBConn
   var dogs []m.Dogs

   db.Find(&dogs) //10ตัว

   type DogsRes struct {
       Name  string `json:"name"`
       DogID int    `json:"dog_id"`
       Type  string `json:"type"`
   }

   var dataResults []DogsRes
   for _, v := range dogs { //1 inet 112 //2 inet1 113
       typeStr := ""
       if v.DogID == 111 {
           typeStr = "red"
       } else if v.DogID == 113 {
           typeStr = "green"
       } else if v.DogID == 999 {
           typeStr = "pink"
       } else {
           typeStr = "no color"
       }

       d := DogsRes{
           Name:  v.Name,  //inet
           DogID: v.DogID, //112
           Type:  typeStr, //no color
       }
       dataResults = append(dataResults, d)
       // sumAmount += v.Amount
   }

   type ResultData struct {
       Data  []DogsRes `json:"data"`
       Name  string    `json:"name"`
       Count int       `json:"count"`
   }
   r := ResultData{
       Data:  dataResults,
       Name:  "golang-test",
       Count: len(dogs), //หาผลรวม,
   }
   return c.Status(200).JSON(r)
}

func Factorial(c *fiber.Ctx) error {
	num, err := strconv.Atoi(c.Params("number"))
	if err != nil || num < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid non-negative integer"})
	}

	fact := 1
	for i := 2; i <= num; i++ {
		fact *= i
	}

	return c.JSON(fiber.Map{
		"number":    num,
		"factorial": fact,
	})
}







