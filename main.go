package main

import (
	"net/http"
	"strconv"

	"./entities"
	rforest "./random_forest"
	"github.com/gin-gonic/gin"
)

// RandomForestTrain , server.GET(train)
func RandomForestTrain(context *gin.Context) {

}

// RandomForestPredict , server.POST(predict)
func RandomForestPredict(context *gin.Context) {

	var pacient entities.Pacient

	if context.Bind(&pacient) == nil {

		url := "https://raw.githubusercontent.com/Malvodio/TF_Funda_Videojuegos/master/diabetes.csv"
		df, _ := rforest.LoadCSV(url)

		feature := []string{
			strconv.Itoa(pacient.Pregnancies),
			strconv.Itoa(pacient.Glucose),
			strconv.Itoa(pacient.BloodPressure),
			strconv.Itoa(pacient.SkinThickness),
			strconv.Itoa(pacient.Insulin),
			strconv.FormatFloat(pacient.BMI, 'f', 6, 64),
			strconv.FormatFloat(pacient.DiabetesPedigreeFunction, 'f', 6, 64),
			strconv.Itoa(pacient.Age),
		}

		inputs, targets, _, _, _, _ := rforest.TrainTestSplit(df, 0.8)

		forest := rforest.BuildForest(inputs, targets, 10, 500, len(inputs[0]))

		X := make([]interface{}, 0)

		for _, x := range feature {
			X = append(X, x)
		}

		result := forest.Predicate(X)

		message := ""

		if result == "0" {
			message = "El paciente podria no padecer diabetes."
		} else {
			message = "El paciente podria padecer diabetes."
		}

		context.JSON(http.StatusOK, gin.H{"result": result, "message": message})
	}

}

func main() {
	server := gin.Default()

	server.GET("/random_forest/train", RandomForestTrain)
	server.POST("/random_forest/predict", RandomForestPredict)

	server.Run()
}
