package golimg

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

/*
colornames: https://cs.opensource.google/go/x/image/+/master:colornames/table.go;l=8
*/

type DrawText struct {
	SrcImgPath      string  `json:"src_img_path"`
	Dpi             float64 `json:"dpi"`
	FontPath        string  `json:"font_path"` // only TTF for now
	Hinting         string  `json:"hinting"`
	FontSize        float64 `json:"font_size"`
	FontSpacing     float64 `json:"font_spacing"`
	FontColorName   string  `json:"font_color_name"` // for names in colorname
	MaxCharsPerLine int     `json:"max_chars_per_line"`
	WhiteOnBlack    bool    `json:"white_on_black"`
	TextTitle       string  `json:"text_title"`
	FontColor       *image.Uniform
}

const (
	DefaultDpi             = 72
	DefaultFontFile        = "fonts/FFF_Tusj.ttf"
	DefaultHinting         = "none" // values: none, full
	DefaultFontSize        = 12     // value in points
	DefaultFontSpacing     = 1.25   // e.g. 2 means double spaced
	DefaultMaxCharsPerLine = 16
	DefaultWhiteOnBlack    = false // useful when creating image as well
	DefaultTextPointX      = 10
)

// Apply default values for DrawText instance if required
func (drawtxt *DrawText) applyDefaultsIfNeeded() {
	if drawtxt.Dpi == 0.0 {
		drawtxt.Dpi = DefaultDpi
	}
	if drawtxt.FontPath == "" {
		drawtxt.FontPath = DefaultFontFile
	}
	if drawtxt.Hinting == "" {
		drawtxt.Hinting = DefaultHinting
	}
	if drawtxt.FontSize == 0.0 {
		drawtxt.FontSize = DefaultFontSize
	}
	if drawtxt.FontSpacing == 0.0 {
		drawtxt.FontSpacing = DefaultFontSpacing
	}
	if drawtxt.MaxCharsPerLine == 0 {
		drawtxt.MaxCharsPerLine = DefaultMaxCharsPerLine
	}
}

// Read the font data.
func (drawtxt *DrawText) GetFont() (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile(drawtxt.FontPath)
	if err != nil {
		return &truetype.Font{}, err
	}
	return truetype.Parse(fontBytes)
}

// Get Foreground Color, for text
func (drawtxt *DrawText) SetFgColor() {
	if drawtxt.FontColorName != "" {
		drawtxt.FontColor = &image.Uniform{C: colornames.Map["yellow"]}
	}
	if drawtxt.WhiteOnBlack {
		drawtxt.FontColor = image.White
	}
	drawtxt.FontColor = image.Black
}

// Get Background Color, for image if created from scratch
func (drawtxt *DrawText) getBgColor() *image.Uniform {
	if drawtxt.WhiteOnBlack {
		return image.Black
	}
	return image.White
}

// Draw the background and the guidelines.
func (drawtxt *DrawText) CreateBgImage(imgW, imgH int) *image.RGBA {
	drawtxt.SetFgColor()
	bg := drawtxt.getBgColor()

	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	if drawtxt.WhiteOnBlack {
		ruler = color.RGBA{0x22, 0x22, 0x22, 0xff}
	}
	rgba := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}
	return rgba
}

// Load an existing Background image
func (drawtxt *DrawText) LoadBgImage(imgpath string) (*image.RGBA, error) {
	if imgpath == "" {
		return drawtxt.CreateBgImage(640, 480), nil
	}
	drawtxt.SetFgColor()

	f, err := os.Open(imgpath)
	if err != nil {
		return &image.RGBA{}, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return &image.RGBA{}, err
	}
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba, nil
	}
	return &image.RGBA{}, errors.New("Failed to convert loaded image to RGBA.")
}

// Get Font Drawer instance with destination image & font configs
func (drawtxt *DrawText) GetFontDrawer(rgba *image.RGBA) (*font.Drawer, error) {
	f, errFont := drawtxt.GetFont()
	if errFont != nil {
		return &font.Drawer{}, errFont
	}

	h := font.HintingNone
	switch drawtxt.Hinting {
	case "full":
		h = font.HintingFull
	}
	d := &font.Drawer{
		Dst: rgba,
		Src: drawtxt.FontColor,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    drawtxt.FontSize,
			DPI:     drawtxt.Dpi,
			Hinting: h,
		}),
	}
	return d, nil
}

func (drawtxt *DrawText) addLine(drawer *font.Drawer, x, y int, txt string) {
	drawer.Dot = fixed.P(x, y)
	drawer.DrawString(txt)
}

// Draw the text.
func (drawtxt *DrawText) AddText(drawer *font.Drawer, max image.Point, text string) error {
	y := 10 + int(math.Ceil(drawtxt.FontSize*drawtxt.Dpi/72))
	dy := int(math.Ceil(drawtxt.FontSize * drawtxt.FontSpacing * drawtxt.Dpi / 72))
	if drawtxt.TextTitle != "" {
		drawer.Dot = fixed.Point26_6{
			X: (fixed.I(max.X) - drawer.MeasureString(drawtxt.TextTitle)) / 2,
			Y: fixed.I(y),
		}
		drawer.DrawString(drawtxt.TextTitle)
		y += dy
	}

	var words = strings.Split(text, " ")
	var s, nextS string
	var currentLineLength int
	var lastIndex = len(words) - 1
	for idx, w := range words {
		currentLineLength += len(w)
		if currentLineLength < drawtxt.MaxCharsPerLine {
			s += " " + strings.TrimSpace(w)
			if idx < lastIndex {
				continue
			}
		} else {
			if idx == lastIndex {
				drawtxt.addLine(drawer, DefaultTextPointX, y, s)
				y += dy
				s = strings.TrimSpace(w)
			}
			nextS = w
		}
		if s == "" { // if there is just one char-list but too long
			s = w
			nextS = ""
		}
		drawtxt.addLine(drawer, DefaultTextPointX, y, s)
		y += dy
		s = nextS
		currentLineLength = 0
	}
	return nil
}

// Save that RGBA image to disk.
func (drawtxt *DrawText) SaveImage(rgba *image.RGBA, outFilepath string) error {
	outFile, err := os.Create(outFilepath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (drawtxt *DrawText) CreateImageWithText(text, saveAs string) error {
	drawtxt.applyDefaultsIfNeeded()

	rgba, errImg := drawtxt.LoadBgImage(drawtxt.SrcImgPath)
	if errImg != nil {
		return errImg
	}

	drawer, errDrawer := drawtxt.GetFontDrawer(rgba)
	if errDrawer != nil {
		return errDrawer
	}

	drawtxt.AddText(drawer, rgba.Rect.Max, text)

	errSave := drawtxt.SaveImage(rgba, saveAs)
	if errSave != nil {
		return errSave
	}
	fmt.Println("Wrote OK:", saveAs)
	return nil
}
