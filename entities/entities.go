package entities

type Pacient struct {
	Pregnancies              int     `json:"pregnancies"`
	Glucose                  int     `json:"glucose"`
	BloodPressure            int     `json:"blood_pressure"`
	SkinThickness            int     `json:"skin_thickness"`
	Insulin                  int     `json:"insulin"`
	BMI                      float64 `json:"bmi"`
	DiabetesPedigreeFunction float64 `json:"diabetes_pedigree_function"`
	Age                      int     `json:"age"`
}
