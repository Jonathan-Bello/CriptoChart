package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Jonathan-Bello/CriptoChart/model"
	"github.com/labstack/echo/v4"
)

// import "github.com/Jonathan-Bello/CriptoChart/model"

type responses []model.Response

func cripto(c echo.Context) error {
	currency := c.Param("currency")
	startDate := c.Param("startdate")

	url := urlNomics + "&ids=" + currency + "&start=" + startDate
	log.Print(url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	cliente := http.Client{}

	resp, err := cliente.Do(req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	dataResponse := responses{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&dataResponse)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	dataChart := createDataChart(dataResponse[0])
	chart := createChart(dataChart)

	// c.JSON(http.StatusOK, z)
	c.HTML(http.StatusOK, chart)
	return nil
}
func createDataChart(resp model.Response) string {
	dataChartHeaders := "['Fecha', '" + resp.Currency + "'],"
	var dataChart string

	for i, time := range resp.Timestamps {
		dataChart += "['" + time.Format("2006-01-02") + "', " + resp.Prices[i] + "],"
	}
	dataChart = strings.TrimSuffix(dataChart, ",")

	dataC := strings.Join([]string{dataChartHeaders, dataChart}, "")

	log.Printf("%v", dataC)

	return dataC
}

func createChart(dataChart string) string {
	htmlChart := `<html>
	<head>
		<script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
		<script type="text/javascript">
			google.charts.load('current', {'packages':['imagelinechart']});
			google.charts.setOnLoadCallback(drawChart);

			function drawChart() {
				var data = google.visualization.arrayToDataTable([` + dataChart + `]);

				var options = {
					title: 'Valor de la criptomoneda(USD)',
					min: 0,
					width: 1000,
					height: 500,
					legend: 'none',
					backgroundColor : '#f2f2f2',
					colors : ['#1092EF']
				};

				var chart = new google.visualization.ImageLineChart(document.getElementById('line_chart'));

				chart.draw(data, options);
			}
		</script>
	</head>
	<body>
		<div id="line_chart" style="width: 1000px; height: 600px"></div>
	</body>
	</html>
	`
	return htmlChart
}
