package main

import (
	"fmt"
	"github.com/wcharczuk/go-chart"
	"math"
	"net/http"
)

func fx(x float64) float64 {
	return math.Cos(2 * x)
}

func approximate(x float64, A []float64) float64 {
	Px := 0.0
	for i := 0; i < len(A); i++ {
		Px += A[i] * math.Pow(x, float64(i))
	}
	return Px
}

func sum(a []float64, j, k float64) float64 {
	sum := 0.0
	for _, val := range a {
		sum += math.Pow(val, j+k)
	}
	return sum
}

var xAxis []float64
var fxAxis []float64
var px []float64

func main() {
	var n = 8 // степень многочлена
	var m = 30

	x := make([]float64, 0)
	f := make([]float64, 0)
	for i := 0.0; i < 2*math.Pi; i += 0.20943951 {
		x = append(x, i)
		f = append(f, fx(i))
	}
	//матрица Грамма
	//g[i][j] = Интеграл от 0 до Пи x^(i+j)
	g := make([][]float64, 0)
	d := make([]float64, 0)

	for k := 0; k < n+1; k++ {
		gSec := make([]float64, n+1)
		for j := 0; j < n+1; j++ {
			gSec[j] = func(a []float64, j, k float64) float64 {
				sum := 0.0
				for _, val := range a {
					sum += math.Pow(val, j+k)
				}
				return sum
			}(x, float64(j), float64(k))

		}
		g = append(g, gSec)
		d = append(d, func(a, f []float64, k float64) float64 {
			sum := 0.0
			for i, _ := range a {
				sum += f[i] * math.Pow(a[i], k)
			}
			return sum
		}(x, f, float64(k)))
	}

	a := Execute(g, d)
	Q := 0.0
	for i := 0; i < m+1; i++ {
		S := 0.0
		for j := 0; j < n+1; j++ {
			S += a[j] * math.Pow(x[i], float64(j))
		}
		Q += math.Pow(f[i]-S, 2)
	}

	for i := 0.0; i < 2*math.Pi; i += 0.01 {
		xAxis = append(xAxis, i)
		fxAxis = append(fxAxis, fx(i))
	}

	for _, val := range xAxis {
		px = append(px, approximate(val, a))
	}

	fmt.Println("близость: ", Q)
	fmt.Println("a: ", a)
	fmt.Println("f(x): ", fxAxis[:5])
	fmt.Println("p(x): ", px[:5])

	http.HandleFunc("/", drawChart)
	http.ListenAndServe(":8000", nil)
}

func drawChart(res http.ResponseWriter, req *http.Request) {

	/*
	   The below will draw the same chart as the `basic` example, except with both the x and y axes turned on.
	   In this case, both the x and y axis ticks are generated automatically, the x and y ranges are established automatically, the canvas "box" is adjusted to fit the space the axes occupy so as not to clip.
	*/

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{
				//Show: true, //enables / displays the x-axis
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				//Show: true, //enables / displays the y-axis
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					//Show:        true,
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					StrokeWidth: 5.0,
					//FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: xAxis,
				YValues: fxAxis,
			},
		},
	}
	res.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, res)
}
