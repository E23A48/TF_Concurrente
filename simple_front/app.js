
$(document).ready(function() {


    $( "#enviar" ).click(function(event) {
        event.preventDefault()

        axios.post('http://localhost:8080/random_forest/predict', {
        "pregnancies": parseInt($("#1").val(), 10),
        "glucose": parseInt($("#2").val(), 10),
        "blood_pressure": parseInt($("#3").val(), 10),
        "skin_thickness": parseInt($("#4").val(), 10),
        "insulin": parseInt($("#5").val(), 10),
        "bmi": parseFloat($("#6").val()),
        "diabetes_pedigree_function": parseFloat($("#7").val()),
        "age": parseInt($("#8").val(), 10)
    })
      .then(function (response) {
        console.log(response.data.message);
        if (response.data.result == 0) {
            Swal.fire(
                'Resultado: ' + response.data.result, 
                response.data.message,
                'success'
              )
        } else {
            Swal.fire(
                'Resultado: ' + response.data.result,
                response.data.message,
                'warning'
            )
        }
      })
      .catch(function (error) {
        Swal.fire(
            'Ocurrio un error',
            'No se pudo precesas tu solicitud',
            'error'
        )
      });
    });

});