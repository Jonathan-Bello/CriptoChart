package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/Jonathan-Bello/CriptoChart/model"
)

type responses []model.Response

// httpRequest is a helper for realize and return a http request
func httpRequest(methond, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(methond, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreateChart is a handler that returns a chart in format image
func CreateChart(c echo.Context) error {
	currency := c.Param("currency")
	startDate := c.Param("startdate") + timeZone
	endDate := c.Param("enddate")

	var url string
	if endDate == "" {
		url = urlNomics + "&ids=" + currency + "&start=" + startDate
	} else {
		url = urlNomics + "&ids=" + currency + "&start=" + startDate + "&end=" + endDate + timeZone
	}

	log.Print(url)

	resp, err := httpRequest(http.MethodGet, url, nil)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	dataResponse := responses{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&dataResponse)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var dataChart string
	if endDate == "" {
		dataChart = createDataChart(dataResponse[0])
	} else {
		dataChart = createDataChart(dataResponse[0])
	}
	chart := htmlChart(dataChart)
	c.HTML(http.StatusOK, chart)
	return nil
}

// createDataChart prepare a string to create the points on the chart
func createDataChart(resp model.Response) string {
	/*
		array that helps create date points to avoid deformation of chart
	*/
	parts := [9]int{
		0,
		int(math.Round(float64(len(resp.Timestamps)) * 1 / 8)),
		int(math.Round(float64(len(resp.Timestamps)) * 2 / 8)),
		int(math.Round(float64(len(resp.Timestamps)) * 3 / 8)),
		int(math.Round(float64(len(resp.Timestamps)) * 4 / 8)),
		int(math.Round(float64(len(resp.Timestamps)) * 5 / 8)),
		int(math.Round(float64(len(resp.Timestamps)) * 6 / 8)),
		int(math.Round(float64(len(resp.Timestamps)) * 7 / 8)),
		len(resp.Timestamps) - 1,
	}
	dataChartHeaders := "['Fecha', '" + resp.Currency + "'],"
	var dataChart string

	for i, time := range resp.Timestamps {
		if containsInt(i, parts[:]) {
			dataChart += "['" + time.Format("2006-01-02") + "', " + resp.Prices[i] + "],"
		}
		dataChart += "['" + time.Format("") + "', " + resp.Prices[i] + "],"
	}

	dataChart = strings.TrimSuffix(dataChart, ",")
	dataC := strings.Join([]string{dataChartHeaders, dataChart}, "")
	return dataC
}

// htmlChart create a HTML struct with the javascript to create the chart picture
func htmlChart(dataChart string) string {
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

// containsInt checks if an array contains a specific number
func containsInt(n int, array []int) bool {
	for _, v := range array {
		if v == n {
			return true
		}
	}
	return false
}
