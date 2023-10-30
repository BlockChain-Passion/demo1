package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Herd struct {
	XMLName xml.Name `xml:"herd"`
	Text    string   `xml:",chardata"`
	Labyak  []struct {
		Text            string  `xml:",chardata"`
		Name            string  `xml:"name,attr"`
		Age             float64 `xml:"age,attr"`
		Sex             string  `xml:"sex,attr"`
		Age_Last_Shaved float64
	} `xml:"labyak"`
}

type Labyak1 struct {
	Name            string
	Age             float64
	Sex             string
	Age_Last_Shaved float64
}

var Labyak1_list []Labyak1

func main() {
	readXML("input.xml")
	router := gin.Default()
	router.GET("/yak-shop/load", resetData)
	router.GET("/yak-shop/stock/:id", getStock)
	router.GET("/yak-shop/herd/:id", getherd)
	router.Run("localhost:8080")

}

func readXML(file_name string) {
	xmlFile, err := os.Open(file_name)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully opened input.xml")

	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var xmlDt []Herd

	err = xml.Unmarshal([]byte(byteValue), &xmlDt)

	fmt.Printf("Read data  ")
	Labyak1_list = []Labyak1{}
	for _, record := range xmlDt {
		for _, Lab := range record.Labyak {
			age := Lab.Age
			Labyak11 := Labyak1{Age: age, Name: Lab.Name, Sex: Lab.Sex, Age_Last_Shaved: 0.00} // Lab.Age_Last_Shaved
			Labyak1_list = append(Labyak1_list, Labyak11)
		}
	}

	fmt.Printf("main method we are here %d", len(Labyak1_list))

}

func resetData(c *gin.Context) {
	fmt.Printf("getAlbums we are here %d", len(Labyak1_list))
	readXML("input.xml")
	c.IndentedJSON(205, Labyak1_list)
}

// *stock
func getStock(c *gin.Context) {
	id := c.Param("id")
	currDay, _ := strconv.Atoi(id)
	c.JSON(http.StatusOK, gin.H{"milk": getmilkstock(currDay), "skins": getSheepSkin(currDay)})
}

func setAge(currentday int) {
	for i, rec := range Labyak1_list {
		rec.Age = rec.Age + float64(currentday)/100
		Labyak1_list[i].Age = rec.Age
	}
}

func getSheepSkin(currentday int) float64 {
	totalSkins := 0
	var nextSkinDay float64
	currDaysInAge := 0.0
	for i, rec := range Labyak1_list {
		initialAgeInDays := rec.Age * 100
		if initialAgeInDays >= 100 {
			totalSkins = totalSkins + 1
		}
		currDaysInAge = initialAgeInDays
		nextSkinDay = 8.00 + float64(currDaysInAge)*0.01

		var currDaysInAge1 float64
		for (float64(currentday)-1)-nextSkinDay > 0 {
			totalSkins = totalSkins + 1
			currDaysInAge1 = float64(currDaysInAge) + float64(nextSkinDay)
			nextSkinDay += (8 + float64(currDaysInAge1)*float64(0.01))
		}
		Labyak1_list[i].Age_Last_Shaved = float64(nextSkinDay)
	}
	return nextSkinDay
}

func getmilkstock(currentday int) float64 {
	var totalMilk float64
	var secondPart float64

	totalMilk = 0.00
	for _, rec := range Labyak1_list {
		initialAgeInDays := rec.Age * 100
		currDaysInAge := initialAgeInDays + float64(currentday)

		if currDaysInAge >= 1000 {
			break
		}
		firstPart := (50 * currentday)
		secondPart = (initialAgeInDays*float64(currentday) + ((float64(currentday) * (float64(currentday) + 1)) / 2))
		secondPart = secondPart * 0.03
		totalMilk = totalMilk + float64(firstPart) - secondPart
	}
	setAge(currentday)

	return totalMilk
}

func getherd(c *gin.Context) {
	id := c.Param("id")
	currDay, _ := strconv.Atoi(id)
	getmilkstock(currDay)
	getSheepSkin(currDay)
	c.IndentedJSON(http.StatusOK, Labyak1_list)
}
