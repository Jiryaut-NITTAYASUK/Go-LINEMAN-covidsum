package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var lists Datas

type Datas struct {
	Data []Data `json:"Data"`
}

type Data struct {
	ConfirmDate    string `json:"ConfirmDate"`
	No             int    `json:"No"`
	Age            *int64 `json:"Age"`
	Gender         string `json:"Gender"`
	GenderEn       string `json:"GenderEn"`
	Nation         string `json:"Nation"`
	NationEn       string `json:"NationEn"`
	Province       string `json:"Province"`
	ProvinceId     int    `json:"ProvinceId"`
	District       string `json:"District"`
	ProvinceEn     string `json:"ProvinceEn"`
	StatQuarantine int    `json:"StatQuarantine"`
}

func greeting(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Greeting"})
}

func summarizeCOVID19(c *gin.Context){
	attended := make(map[string]int)
	for i:=0 ; i<len(lists.Data); i++ {
		a := lists.Data[i].ProvinceEn
		if lists.Data[i].ProvinceEn == "" {
			a = "N/A"
		}
		if attended[a] == 0 {
            attended[a] = 1
        } else {
            attended[a] = attended[a] + 1
        } 
    }
	
	attended2 := map[string]int {
        "0-30": 0,
        "31-60": 0,
        "61+": 0,
        "N/A": 0,
    }
    for i:=0 ; i<len(lists.Data); i++ {
        a := lists.Data[i].Age
        if a != nil{
			if *a>=0 && *a<=30 {
				attended2["0-30"]+=1
			} else if *a>=31 && *a<=60 {
				attended2["31-60"]+=1
			} else if *a >= 61 {
				attended2["61+"]+=1
			}
		}else{
			attended2["N/A"]+=1
		}
    }
	result := map[string]interface{}{
		"Province": attended,
		"AgeGroup": attended2,
	}

	c.JSON(http.StatusOK, result)
}

func main() {
	jsonFile, err := http.Get("http://static.wongnai.com/devinterview/covid-cases.json")
	if err != nil {
	   log.Fatalln(err)
	}
	jsonFile.Header = http.Header{
		"Content-Type": {"application/json"},
	}

	body, err := io.ReadAll(jsonFile.Body)
	if err != nil {
	   log.Fatalln(err)
	}

	json.Unmarshal(body, &lists)

	r := gin.Default()
	r.GET("/", greeting)
	r.GET("/covid/summary", summarizeCOVID19)
	r.Run("localhost:8080")
}