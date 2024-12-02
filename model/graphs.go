package model

import (
	"fmt"
	"math/rand/v2"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func TestBarChart() *charts.Bar {
	//create bar chart
	bar := charts.NewBar()

	//options and title
	bar.SetGlobalOptions(
	charts.WithTitleOpts(opts.Title{
		Title: "Example title",
		TitleStyle: &opts.TextStyle{
			Color: opts.RGBColor(255, 255, 255),
		},
		Subtitle: "Example subtitle",
		SubtitleStyle: &opts.TextStyle{
			Color: opts.RGBColor(255, 255, 255),
		},
		BorderColor: opts.RGBColor(255, 255, 255),
	}))

	bar.SetGlobalOptions(
	charts.WithXAxisOpts(opts.XAxis{
		AxisLabel: &opts.AxisLabel{
			Color: opts.RGBColor(255, 255, 255),
		},
	}))

	bar.SetGlobalOptions(
		charts.WithYAxisOpts(
			opts.YAxis{
				AxisLabel: &opts.AxisLabel{
					Color: opts.RGBColor(255, 255, 255),
				},
			},
		),
	)

	bar.SetGlobalOptions(
		charts.WithLegendOpts(
			opts.Legend{
				TextStyle: &opts.TextStyle{
					Color: opts.RGBColor(255, 255, 255),
				},
			},
		),
	)

	//add data
	data := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		data = append(data, opts.BarData{Value: rand.IntN(300)})
	}

	bar.SetXAxis([]string{"xaxis1", "xaxis2", "xaxis3", "axis4", "axis5", "axis6", "axis7", "axis8"}).
		AddSeries("Series1", data)

	return bar
}

func StyledBarChart() *charts.Bar{
	//create bar chart
	bar := charts.NewBar()


	bar.SetGlobalOptions(
	charts.WithXAxisOpts(opts.XAxis{
		AxisLabel: &opts.AxisLabel{
			Color: opts.RGBColor(255, 255, 255),
		},
	}))

	bar.SetGlobalOptions(
		charts.WithYAxisOpts(
			opts.YAxis{
				AxisLabel: &opts.AxisLabel{
					Color: opts.RGBColor(255, 255, 255),
				},
			},
		),
	)

	bar.SetGlobalOptions(
		charts.WithLegendOpts(
			opts.Legend{
				TextStyle: &opts.TextStyle{
					Color: opts.RGBColor(255, 255, 255),
				},
			},
		),
	)
	return bar
}

func SalesFromDateBarChart(data []MostSoldProductData) *charts.Bar {
	fmt.Printf("%v \n \n", data[0])

	bar := StyledBarChart()

	chartData := make([]opts.BarData, 0)
	chartTitle := data[0].Name
	xAxisLabels := []string{}


	for i := range(data) {
		xAxisLabels = append(xAxisLabels, data[i].Date)
		chartData = append(chartData, opts.BarData{Value: data[i].Quantity})
	}

	bar.SetXAxis(xAxisLabels).
		AddSeries(chartTitle, chartData)

	bar.SetGlobalOptions(
		charts.WithTitleOpts(
			opts.Title{
				Title: "Cantidad de ventas",
				TitleStyle: &opts.TextStyle{
					Color: opts.RGBColor(255, 255, 255),
				},
				Subtitle: chartTitle,
				SubtitleStyle: &opts.TextStyle{
					Color: opts.RGBColor(255, 255, 255),
				},
			},
		),
	)
	return bar
}

func ProductDeltaTimespanChart(data []ProductSaleDelta) *charts.Bar {
	fmt.Printf("%v \n \n", data[0])

	bar := StyledBarChart()

	salesData := make([]opts.BarData, 0)
	variationData := make([]opts.BarData,0)
	chartTitle := data[0].Name
	xAxisLabels := []string{}


	for i := range(data) {
		xAxisLabels = append(xAxisLabels, data[i].Date)
		salesData = append(salesData, opts.BarData{Value: data[i].CurrentSales, Name: "Ventas"})
		variationData = append(variationData, opts.BarData{Value: data[i].MonthlyVariance, Name: "Diferencia mensual"})
	}

	bar.SetXAxis(xAxisLabels).
		AddSeries("Ventas", salesData).
		AddSeries("Diferencia del periodo", variationData)

	bar.SetGlobalOptions(
		charts.WithTitleOpts(
			opts.Title{
				Title: "Diferencia en ventas de producto",
				TitleStyle: &opts.TextStyle{
					Color: opts.RGBColor(255, 255, 255),
				},
				Subtitle: chartTitle,
				SubtitleStyle: &opts.TextStyle{
					Color: opts.RGBColor(255, 255, 255),
				},
			},
		),
	)
	return bar
}

func ProfitDeltaTimespanChart(data []ProfitDelta) *charts.Bar {
	fmt.Printf("%v \n \n", data[0])

	bar := StyledBarChart()

	salesData := make([]opts.BarData, 0)
	variationData := make([]opts.BarData,0)
	xAxisLabels := []string{}

	for i := range(data) {
		xAxisLabels = append(xAxisLabels, data[i].Date)
		salesData = append(salesData, opts.BarData{Value: data[i].Profits, Name: "Ingresos"})
		variationData = append(variationData, opts.BarData{Value: data[i].DifferencePrevious, Name: "Diferencia de la fecha anterior"})
	}

	bar.SetXAxis(xAxisLabels).
		AddSeries("Ingresos", salesData).
		AddSeries("Diferencia de la fecha", variationData)

	bar.SetGlobalOptions(
		charts.WithTitleOpts(
			opts.Title{
				Title: "Diferencia de ingresos",
				TitleStyle: &opts.TextStyle{
					Color: opts.RGBColor(255, 255, 255),
				},
			},
		),
	)
	return bar
}

func SalesDeltaTimespanChart(data []ProductSaleDelta) *charts.Bar {
	fmt.Printf("%v \n \n", data[0])

	bar := StyledBarChart()

	subtitle := data[0].Name
	salesData := make([]opts.BarData, 0)
	variationData := make([]opts.BarData,0)
	xAxisLabels := []string{}


	for i := range(data) {
		xAxisLabels = append(xAxisLabels, data[i].Date)
		salesData = append(salesData, opts.BarData{Value: data[i].CurrentSales, Name: "Ingresos"})
		variationData = append(variationData, opts.BarData{Value: data[i].MonthlyVariance, Name: "Diferencia de la fecha anterior"})
	}

	bar.SetXAxis(xAxisLabels).
		AddSeries("Ingresos", salesData).
		AddSeries("Diferencia de la fecha", variationData)

	bar.SetGlobalOptions(
		charts.WithTitleOpts(
			opts.Title{
				Title: "Diferencia en ventas",
				TitleStyle: &opts.TextStyle{
					Color: opts.RGBColor(255, 255, 255),
				},
				Subtitle: subtitle,
				SubtitleStyle: &opts.TextStyle{
					Color: opts.RGBColor(255, 255, 255),
				},
			},
		),
	)
	return bar
}

func MoneySpentProductDeltaTimespanChart(data []MoneySpentProductDelta) *charts.Bar {
	fmt.Printf("%v \n \n", data[0])

	bar := StyledBarChart()

	subtitle := data[0].Name
	salesData := make([]opts.BarData, 0)
	variationData := make([]opts.BarData,0)
	xAxisLabels := []string{}


	for i := range(data) {
		xAxisLabels = append(xAxisLabels, data[i].Date)
		salesData = append(salesData, opts.BarData{Value: data[i].Expenses, Name: "Gastos"})
		variationData = append(variationData, opts.BarData{Value: data[i].Difference, Name: "Diferencia de la fecha anterior"})
	}

	bar.SetXAxis(xAxisLabels).
		AddSeries("Gastos", salesData).
		AddSeries("Diferencia de la fecha", variationData)

	bar.SetGlobalOptions(
		charts.WithTitleOpts(
			opts.Title{
				Title: "Diferencia en gastos del producto",
				TitleStyle: &opts.TextStyle{
					Color: opts.RGBColor(255, 255, 255),
				},
				Subtitle: subtitle,
				SubtitleStyle: &opts.TextStyle{
					Color: opts.RGBColor(255, 255, 255),
				},
			},
		),
	)
	return bar
}

func MoneySpentDeltaTimespanChart(data []MoneySpentDelta) *charts.Bar {
	fmt.Printf("%v \n \n", data[0])

	bar := StyledBarChart()

	salesData := make([]opts.BarData, 0)
	variationData := make([]opts.BarData,0)
	xAxisLabels := []string{}


	for i := range(data) {
		xAxisLabels = append(xAxisLabels, data[i].Date)
		salesData = append(salesData, opts.BarData{Value: data[i].Sales, Name: "Gastos"})
		variationData = append(variationData, opts.BarData{Value: data[i].Difference, Name: "Diferencia de la fecha anterior"})
	}

	bar.SetXAxis(xAxisLabels).
		AddSeries("Gastos", salesData).
		AddSeries("Diferencia de la fecha", variationData)

	bar.SetGlobalOptions(
		charts.WithTitleOpts(
			opts.Title{
				Title: "Diferencia en gastos",
				TitleStyle: &opts.TextStyle{
					Color: opts.RGBColor(255, 255, 255),
				},
				Subtitle: "Gastos totales",
				SubtitleStyle: &opts.TextStyle{
					Color: opts.RGBColor(255, 255, 255),
				},
			},
		),
	)
	return bar
}
