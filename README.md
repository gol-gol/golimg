
## golimg

> few utility functions for specific Image Manipulation in Go

### Public Functions

> `DrawText{SrcImgPath, Dpi, FontPath, Hinting, FontSize, FontSpacing, FontColorName, MaxCharsPerLine, WhiteOnBlack, TextTitle}`

* `DrawText.CreateImageWithText(text, saveAs string)` is the primary function allowing creating Image with Text overlay

> it will create a default 640x480 image if no `SrcImgPath` provided, to create custom image size.. other public functions can be used in combination which are chained together in here

* `DrawText.GetFont() (*truetype.Font, error)`
* `DrawText.SetFgColor()`
* `DrawText.CreateBgImage(imgW, imgH int) *image.RGBA`
* `DrawText.LoadBgImage(imgpath string) (*image.RGBA, error)`
* `DrawText.GetFontDrawer(rgba *image.RGBA) (*font.Drawer, error)`
* `DrawText.AddText(drawer *font.Drawer, max image.Point, text string) error`
* `DrawText.SaveImage(rgba *image.RGBA, outFilepath string) error`

---
