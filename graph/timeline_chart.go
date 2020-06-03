package graph

import (
	"errors"
	"fmt"
	"io"
	"math"

	chart "github.com/wcharczuk/go-chart"

	"github.com/golang/freetype/truetype"
	"github.com/wcharczuk/go-chart/seq"
	util "github.com/wcharczuk/go-chart/util"
)

type BarSection struct {
	StartValue float64
	Length     float64
}

// BarSections is an array of Value.
type BarSections []BarSection

//// BarSections returns the values.
//func (vs BarSections) Values() []float64 {
//	values := make([]float64, len(vs))
//	for index, v := range vs {
//		values[index] = v.Length
//	}
//
//	return values
//}

//// ValuesNormalized returns normalized values.
//func (vs BarSections) ValuesNormalized() []float64 {
//	return util.Math.Normalize(vs.Values()...)
//}

// Normalize returns the values normalized.
func (vs BarSections) Normalize(minValue, maxValue float64) []BarSection {
	var output []BarSection
	for _, v := range vs {
		if v.Length > 0 {
			output = append(output, BarSection{
				StartValue: util.Math.RoundUp((v.StartValue-minValue)/(maxValue-minValue), 0.0001),
				Length:     util.Math.RoundUp(v.Length/(maxValue-minValue), 0.0001),
			})
		}
	}

	return output
}

// Bar is a bar within a TimelineChart.
type Bar struct {
	Name     string
	Style    chart.Style
	Sections BarSections
}

// GetHeight returns the width of the bar.
func (sb Bar) GetHeight() int {
	return 15
}

// TimelineChart is a chart that draws sections of a bar based on percentages.
type TimelineChart struct {
	Title      string
	TitleStyle chart.Style

	ColorPalette chart.ColorPalette

	Width  int
	Height int
	DPI    float64

	Background chart.Style
	Canvas     chart.Style

	XAxis chart.Style
	YAxis chart.Style

	BarSpacing int

	Font        *truetype.Font
	defaultFont *truetype.Font

	Bars     []Bar
	Elements []chart.Renderable
}

// GetDPI returns the dpi for the chart.
func (sbc TimelineChart) GetDPI(defaults ...float64) float64 {
	if sbc.DPI == 0 {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return chart.DefaultDPI
	}
	return sbc.DPI
}

// GetFont returns the text font.
func (sbc TimelineChart) GetFont() *truetype.Font {
	if sbc.Font == nil {
		return sbc.defaultFont
	}
	return sbc.Font
}

// GetHeight returns the chart width or the default value.
func (sbc TimelineChart) GetWidth() int {
	if sbc.Width == 0 {
		return chart.DefaultChartWidth
	}
	return sbc.Width
}

// GetHeight returns the chart height or the default value.
func (sbc TimelineChart) GetHeight() int {
	if sbc.Height == 0 {
		return chart.DefaultChartHeight
	}
	return sbc.Height
}

// GetBarSpacing returns the spacing between bars.
func (sbc TimelineChart) GetBarSpacing() int {
	if sbc.BarSpacing == 0 {
		return 100
	}
	return sbc.BarSpacing
}

// Render renders the chart with the given renderer to the given io.Writer.
func (sbc TimelineChart) Render(rp chart.RendererProvider, w io.Writer) error {
	if len(sbc.Bars) == 0 {
		return errors.New("please provide at least one bar")
	}

	r, err := rp(sbc.GetWidth(), sbc.GetHeight())
	if err != nil {
		return err
	}

	if sbc.Font == nil {
		defaultFont, err := chart.GetDefaultFont()
		if err != nil {
			return err
		}
		sbc.defaultFont = defaultFont
	}
	r.SetDPI(sbc.GetDPI(chart.DefaultDPI))

	canvasBox := sbc.getAdjustedCanvasBox(r, sbc.getDefaultCanvasBox())
	sbc.drawCanvas(r, canvasBox)
	sbc.drawBars(r, canvasBox)
	sbc.drawXAxis(r, canvasBox)
	sbc.drawYAxis(r, canvasBox)

	sbc.drawTitle(r)
	for _, a := range sbc.Elements {
		a(r, canvasBox, sbc.styleDefaultsElements())
	}

	return r.Save(w)
}

func (sbc TimelineChart) drawCanvas(r chart.Renderer, canvasBox chart.Box) {
	chart.Draw.Box(r, canvasBox, sbc.getCanvasStyle())
}

func (sbc TimelineChart) drawBars(r chart.Renderer, canvasBox chart.Box) {
	yoffset := canvasBox.Bottom
	for _, bar := range sbc.Bars {
		sbc.drawBar(r, canvasBox, sbc.getMinValue(), sbc.getMaxValue(), yoffset, bar)
		yoffset -= sbc.GetBarSpacing() - bar.GetHeight()
	}
}

func (sbc TimelineChart) getMaxValue() float64 {
	maxValue := float64(0)
	for _, bar := range sbc.Bars {
		for _, section := range bar.Sections {
			maxValue = math.Max(maxValue, section.Length+section.StartValue)
		}
	}

	return maxValue
}

func (sbc TimelineChart) getMinValue() float64 {
	maxValue := math.MaxFloat64
	for _, bar := range sbc.Bars {
		for _, section := range bar.Sections {
			maxValue = math.Min(maxValue, section.StartValue)
		}
	}

	return maxValue
}

func (sbc TimelineChart) drawBar(r chart.Renderer, canvasBox chart.Box, minValue, maxValue float64, yoffset int, bar Bar) int {
	barSpacing2 := sbc.GetBarSpacing() >> 1
	bxl := yoffset + barSpacing2
	bxr := bxl + bar.GetHeight()

	normalizedBarComponents := bar.Sections.Normalize(minValue, maxValue)
	for index, bv := range normalizedBarComponents {
		barWidth := int(math.Ceil(bv.Length * float64(canvasBox.Width())))
		barStart := int(math.Ceil(bv.StartValue * float64(canvasBox.Width())))
		barBox := chart.Box{
			Top:    bxl,
			Left:   util.Math.MinInt(barStart, canvasBox.Right-chart.DefaultStrokeWidth),
			Right:  util.Math.MinInt(barStart+barWidth, canvasBox.Right-chart.DefaultStrokeWidth),
			Bottom: bxr,
		}

		chart.Draw.Box(r, barBox, bar.Style.InheritFrom(sbc.styleDefaultsStackedBarValue(index)))
	}

	return bxr
}

func (sbc TimelineChart) drawXAxis(r chart.Renderer, canvasBox chart.Box) {
	if sbc.XAxis.Show {
		axisStyle := sbc.XAxis.InheritFrom(sbc.styleDefaultsAxes())
		axisStyle.WriteToRenderer(r)

		r.MoveTo(canvasBox.Left, canvasBox.Bottom)
		r.LineTo(canvasBox.Right, canvasBox.Bottom)
		r.Stroke()

		r.MoveTo(canvasBox.Left, canvasBox.Bottom)
		r.LineTo(canvasBox.Left, canvasBox.Bottom+chart.DefaultVerticalTickHeight)
		r.Stroke()

		cursor := canvasBox.Left
		for _, bar := range sbc.Bars {

			barLabelBox := chart.Box{
				Top:    canvasBox.Bottom + chart.DefaultXAxisMargin,
				Left:   cursor,
				Right:  cursor + bar.GetHeight() + sbc.GetBarSpacing(),
				Bottom: sbc.GetHeight(),
			}
			if len(bar.Name) > 0 {
				chart.Draw.TextWithin(r, bar.Name, barLabelBox, axisStyle)
			}
			axisStyle.WriteToRenderer(r)
			r.MoveTo(barLabelBox.Right, canvasBox.Bottom)
			r.LineTo(barLabelBox.Right, canvasBox.Bottom+chart.DefaultVerticalTickHeight)
			r.Stroke()
			cursor += bar.GetHeight() + sbc.GetBarSpacing()
		}
	}
}

func (sbc TimelineChart) drawXAxis(r chart.Renderer, canvasBox chart.Box) {
	if sbc.XAxis.Show {
		axisStyle := sbc.XAxis.InheritFrom(sbc.styleDefaultsAxes())
		axisStyle.WriteToRenderer(r)
		r.MoveTo(canvasBox.Right, canvasBox.Top)
		r.LineTo(canvasBox.Right, canvasBox.Bottom)
		r.Stroke()

		r.MoveTo(canvasBox.Right, canvasBox.Bottom)
		r.LineTo(canvasBox.Right+chart.DefaultHorizontalTickWidth, canvasBox.Bottom)
		r.Stroke()

		ticks := seq.RangeWithStep(0.0, 1.0, 0.2)
		for _, t := range ticks {
			axisStyle.GetStrokeOptions().WriteToRenderer(r)
			ty := canvasBox.Bottom - int(t*float64(canvasBox.Height()))
			r.MoveTo(canvasBox.Right, ty)
			r.LineTo(canvasBox.Right+chart.DefaultHorizontalTickWidth, ty)
			r.Stroke()

			axisStyle.GetTextOptions().WriteToRenderer(r)
			text := fmt.Sprintf("%0.0f%%", t*100)

			tb := r.MeasureText(text)
			chart.Draw.Text(r, text, canvasBox.Right+chart.DefaultYAxisMargin+5, ty+(tb.Height()>>1), axisStyle)
		}

	}
}

func (sbc TimelineChart) drawTitle(r chart.Renderer) {
	if len(sbc.Title) > 0 && sbc.TitleStyle.Show {
		r.SetFont(sbc.TitleStyle.GetFont(sbc.GetFont()))
		r.SetFontColor(sbc.TitleStyle.GetFontColor(sbc.GetColorPalette().TextColor()))
		titleFontSize := sbc.TitleStyle.GetFontSize(chart.DefaultTitleFontSize)
		r.SetFontSize(titleFontSize)

		textBox := r.MeasureText(sbc.Title)

		textWidth := textBox.Width()
		textHeight := textBox.Height()

		titleX := (sbc.GetWidth() >> 1) - (textWidth >> 1)
		titleY := sbc.TitleStyle.Padding.GetTop(chart.DefaultTitleTop) + textHeight

		r.Text(sbc.Title, titleX, titleY)
	}
}

func (sbc TimelineChart) getCanvasStyle() chart.Style {
	return sbc.Canvas.InheritFrom(sbc.styleDefaultsCanvas())
}

func (sbc TimelineChart) styleDefaultsCanvas() chart.Style {
	return chart.Style{
		FillColor:   sbc.GetColorPalette().CanvasColor(),
		StrokeColor: sbc.GetColorPalette().CanvasStrokeColor(),
		StrokeWidth: chart.DefaultCanvasStrokeWidth,
	}
}

// GetColorPalette returns the color palette for the chart.
func (sbc TimelineChart) GetColorPalette() chart.ColorPalette {
	if sbc.ColorPalette != nil {
		return sbc.ColorPalette
	}
	return chart.AlternateColorPalette
}

func (sbc TimelineChart) getDefaultCanvasBox() chart.Box {
	return sbc.Box()
}

func (sbc TimelineChart) getAdjustedCanvasBox(r chart.Renderer, canvasBox chart.Box) chart.Box {
	var totalHeight int
	for _, bar := range sbc.Bars {
		totalHeight += bar.GetHeight() + sbc.GetBarSpacing()
	}
	//
	//if sbc.YAxis.Show {
	//	yaxisHeight := chart.DefaultHorizontalTickWidth
	//
	//	axisStyle := sbc.YAxis.InheritFrom(sbc.styleDefaultsAxes())
	//	axisStyle.WriteToRenderer(r)
	//
	//	cursor := canvasBox.Bottom
	//	for _, bar := range sbc.Bars {
	//		if len(bar.Name) > 0 {
	//			barLabelBox := chart.Box{
	//				Top:    cursor + bar.GetHeight() + sbc.GetBarSpacing(), //canvasBox.Bottom + chart.DefaultXAxisMargin,
	//				Left:   canvasBox.Left, //cursor,
	//				Right:  cursor + bar.GetHeight() + sbc.GetBarSpacing(),
	//				Bottom: sbc.GetHeight(),
	//			}
	//			lines := chart.Text.WrapFit(r, bar.Name, barLabelBox.Width(), axisStyle)
	//			linesBox := chart.Text.MeasureLines(r, lines, axisStyle)
	//
	//			yaxisHeight = util.Math.MaxInt(linesBox.Height()+(2*chart.DefaultXAxisMargin), yaxisHeight)
	//		}
	//	}
	//	return chart.Box{
	//		Top:    canvasBox.Top,
	//		Left:   canvasBox.Left,
	//		Right:  canvasBox.Left + totalHeight,
	//		Bottom: sbc.GetHeight() - yaxisHeight,
	//	}
	//}
	return chart.Box{
		Top:    canvasBox.Top,
		Left:   canvasBox.Left,
		Right:  canvasBox.Left + totalHeight,
		Bottom: canvasBox.Bottom,
	}

}

// Box returns the chart bounds as a box.
func (sbc TimelineChart) Box() chart.Box {
	dpr := sbc.Background.Padding.GetRight(10)
	dpb := sbc.Background.Padding.GetBottom(50)

	return chart.Box{
		Top:    sbc.Background.Padding.GetTop(20),
		Left:   sbc.Background.Padding.GetLeft(20),
		Right:  sbc.GetWidth() - dpr,
		Bottom: sbc.GetHeight() - dpb,
	}
}

func (sbc TimelineChart) styleDefaultsStackedBarValue(index int) chart.Style {
	return chart.Style{
		StrokeColor: sbc.GetColorPalette().GetSeriesColor(index),
		StrokeWidth: 3.0,
		FillColor:   sbc.GetColorPalette().GetSeriesColor(index),
	}
}

func (sbc TimelineChart) styleDefaultsTitle() chart.Style {
	return sbc.TitleStyle.InheritFrom(chart.Style{
		FontColor:           chart.DefaultTextColor,
		Font:                sbc.GetFont(),
		FontSize:            sbc.getTitleFontSize(),
		TextHorizontalAlign: chart.TextHorizontalAlignCenter,
		TextVerticalAlign:   chart.TextVerticalAlignTop,
		TextWrap:            chart.TextWrapWord,
	})
}

func (sbc TimelineChart) getTitleFontSize() float64 {
	effectiveDimension := util.Math.MinInt(sbc.GetWidth(), sbc.GetHeight())
	if effectiveDimension >= 2048 {
		return 48
	} else if effectiveDimension >= 1024 {
		return 24
	} else if effectiveDimension >= 512 {
		return 18
	} else if effectiveDimension >= 256 {
		return 12
	}
	return 10
}

func (sbc TimelineChart) styleDefaultsAxes() chart.Style {
	return chart.Style{
		StrokeColor:         chart.DefaultAxisColor,
		Font:                sbc.GetFont(),
		FontSize:            chart.DefaultAxisFontSize,
		FontColor:           chart.DefaultAxisColor,
		TextHorizontalAlign: chart.TextHorizontalAlignCenter,
		TextVerticalAlign:   chart.TextVerticalAlignTop,
		TextWrap:            chart.TextWrapWord,
	}
}
func (sbc TimelineChart) styleDefaultsElements() chart.Style {
	return chart.Style{
		Font: sbc.GetFont(),
	}
}
