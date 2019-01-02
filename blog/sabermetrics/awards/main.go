package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	log.Println(createCorrelationPlot(os.Args[1]))
}

type xy struct {
	x []float64
	y []float64
}

func createCorrelationPlot(csvPath string) error {

	// For every stat in a passed csv file, create
	// a scatter plot showing the correlation between
	// that stat and the votes the players received

	// Open CSV
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Extract data
	data := csv.NewReader(f)
	rows, err := data.ReadAll()
	if err != nil {
		return err
	}

	// Cut headers
	rows = rows[1:]

	// These are the indexes in the CSV data
	// that we care about
	indiciesToStats := map[int]string{
		6:  "WAR",
		7:  "Wins",
		24: "Strikeouts",
		29: "WHIP",
		30: "ERA+",
	}

	// Go through each stat we care about, and create
	// a chart for it (by index is csv data)
	for index, statName := range indiciesToStats {
		// Create a new plot
		p, err := plot.New()
		if err != nil {
			return err
		}

		// Set plot cosmetics
		p.BackgroundColor = color.Gray{Y: 230}
		p.X.Label.Text = "Voting Points"
		p.Y.Label.Text = statName
		p.Title.Text = fmt.Sprintf("Correlation Between %s and Cy Young Votes - 2018", statName)
		p.Title.Font.Size = 25.0
		p.X.Width = 2.0
		p.X.Label.Font.Size = 25.0
		p.X.Tick.LineStyle.Width = 2.0
		p.X.Tick.Label.Font.Size = 16.0
		p.Y.Width = 2.0
		p.Y.Label.Font.Size = 25.0
		p.Y.Tick.LineStyle.Width = 2.0
		p.Y.Tick.Label.Font.Size = 16.0

		data := xy{
			x: make([]float64, len(rows)),
			y: make([]float64, len(rows)),
		}

		// Traverse data, adding points to scatter plot, and regression
		points := make(plotter.XYs, len(rows))
		for i, row := range rows {

			votesFloat, err := strconv.ParseFloat(row[3], 64)
			if err != nil {
				return err
			}

			dataFloat, err := strconv.ParseFloat(row[index], 64)
			if err != nil {
				return err
			}

			points[i].X = votesFloat
			points[i].Y = dataFloat
			data.x[i] = votesFloat
			data.y[i] = dataFloat

		}

		m, c := stat.LinearRegression(data.x, data.y, nil, false)
		line := plotter.NewFunction(func(x float64) float64 { return c*x + m })
		line.Width = 2.0
		line.Color = color.RGBA{R: 255, A: 255}
		line.Dashes = []vg.Length{vg.Points(20), vg.Points(10)}

		fmt.Println(statName, csvPath)

		s, err := plotter.NewScatter(points)
		if err != nil {
			return err
		}

		var weights []float64
		fmt.Println(stat.RSquared(data.x, data.y, weights, m, c))
		fmt.Println()

		s.GlyphStyle.Shape = draw.CircleGlyph{}
		s.GlyphStyle.Radius = 3.0

		p.Add(s, line)

		// Expand each of the axis a little bit
		p.X.Max = p.X.Max + (p.X.Max * 0.1)
		p.X.Min = p.X.Min - (p.X.Min * 0.1)
		p.Y.Max = p.Y.Max + (p.Y.Max * 0.1)
		p.Y.Min = p.Y.Min - (p.Y.Min * 0.1)

		p.Save(10*vg.Inch, 10*vg.Inch, fmt.Sprintf("./graphs/cy-young/%s-%s.png", csvPath, statName))
	}

	return nil
}
