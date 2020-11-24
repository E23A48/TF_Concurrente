package main

import (
	"fmt"
	"net/http"
	"strconv"

	"./consenso"
	"./entities"
	rforest "./random_forest"
	"github.com/gin-gonic/gin"
)

func policyAPIcors() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		context.Next()
	}
}

// RandomForestTrain , server.GET(train)
func RandomForestTrain(context *gin.Context) {

	n_tree, _ := strconv.Atoi(context.Param("n_tree"))

	url := "https://raw.githubusercontent.com/Malvodio/TF_Funda_Videojuegos/master/diabetes.csv"
	df, _ := rforest.LoadCSV(url)

	inputs, targets, _, _, X_test, y_test := rforest.TrainTestSplit(df, 0.8)

	fmt.Println(n_tree)
	forest := rforest.BuildForest(inputs, targets, n_tree, 500, len(inputs[0]))

	test_predicted := make([]string, 0)

	tp, tn, fp, fn := 0.0, 0.0, 0.0, 0.0

	for i := 0; i < len(X_test); i++ {
		fmt.Println(X_test[i])
		output := forest.Predicate(X_test[i])
		expect := y_test[i]

		test_predicted = append(test_predicted, output)

		fmt.Println("Output: ", forest.Predicate(X_test[i]), "Expected: ", expect)
		if expect == "1" {
			if output == expect {
				tp += 1
			} else {
				fp += 1
			}
		} else {
			if output == expect {
				tn += 1
			} else {
				fn += 1
			}
		}
	}

	fmt.Println(tp, tn, fp, fn)

	context.JSON(http.StatusOK, gin.H{
		"res":       "Modelo entrenado correctamente",
		"accuracy":  (tp + tn) / float64(len(y_test)),
		"precision": tp / (tp + fp),
		"recall":    tp / (tp + fn),
		"F1":        2 * ((tp / (tp + fn)) * (tp / (tp + fp))) / ((tp / (tp + fn)) + (tp / (tp + fp))),
	})
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

		_, _, X_train, y_train, _, _ := rforest.TrainTestSplit(df, 0.15)

		forest := rforest.BuildForest(X_train, y_train, 10, 500, len(X_train[0]))

		X := make([]interface{}, 0)

		for _, x := range feature {
			X = append(X, x)
		}

		var result string = forest.Predicate(X)
		result_int, _ := strconv.Atoi(result)

		//----------------------------------------------

		msg := consenso.Tmsg{consenso.Cnum, consenso.LocalAddr, result_int, pacient}
		for _, addr := range consenso.Addrs {
			consenso.Send(addr, msg)
		}

		for consenso.Prediction == -1 {

		}

		resultpred := consenso.Prediction
		consenso.Prediction = -1

		//---------------------------------------------

		message := ""

		if resultpred == 0 {
			message = "El paciente podria no padecer diabetes."
		} else {
			message = "El paciente podria padecer diabetes."
		}

		context.JSON(http.StatusOK, gin.H{"result": resultpred, "message": message})
	}

}

func main() {

	go consenso.GoSV()

	router := gin.Default()

	router.Use(policyAPIcors())

	router.GET("/random_forest/train/:n_tree", RandomForestTrain)
	router.POST("/random_forest/predict", RandomForestPredict)

	router.Run()

}
