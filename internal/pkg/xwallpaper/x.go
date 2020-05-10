package xwallpaper

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xinerama"
	"github.com/nfnt/resize"
)

type XWallpaper struct {
	conn *xgbutil.XUtil
}

func New() *XWallpaper {
	x11 := new(XWallpaper)

	conn, err := xgbutil.NewConn()
	if err != nil {
		log.Fatalf("could not initialise X connection: %v", err)
	}
	x11.conn = conn

	return x11
}

func (xw *XWallpaper) rootRect() image.Rectangle {
	screen := xw.conn.Screen()
	w, h := screen.WidthInPixels, screen.HeightInPixels
	return image.Rect(0, 0, int(w), int(h))
}

func (xw *XWallpaper) headsRect() []image.Rectangle {
	if _, ok := xw.conn.Conn().Extensions["XINERAMA"]; !ok {
		return []image.Rectangle{xw.rootRect()}
	}

	heads, err := xinerama.PhysicalHeads(xw.conn)
	if err != nil {
		log.Fatalf("could not request heads: %v", err)
	}

	rects := make([]image.Rectangle, 0, len(heads))

	for _, head := range heads {
		headX, headY, headW, headH := head.Pieces()
		r := image.Rect(headX, headY, headX+headW, headY+headH)
		rects = append(rects, r)
	}

	return rects
}

func resizeImage(src image.Image, w, h int) image.Image {
	if w == src.Bounds().Dx() || h == src.Bounds().Dy() {
		log.Println("skipping resize")
		return src
	}
	return resize.Resize(uint(w), uint(h), src, resize.Lanczos3)
}

func (xw *XWallpaper) Set(img image.Image) error {

	rootImg := image.NewRGBA(xw.rootRect())

	imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()

	for _, head := range xw.headsRect() {

		headW, headH := head.Dx(), head.Dy()

		var newW, newH int

		switch {
		case imgW*headH > headW*imgH: // image ratio is greater than the head
			newH = headH
		case imgW*headH <= headW*imgH: // image ratio is less or equal to the head
			newW = headW
		}

		resized := resizeImage(img, newW, newH)

		var newSp image.Point
		if newW == 0 {
			newSp = image.Pt((resized.Bounds().Dx()-headW)/2, 0)
		} else {
			newSp = image.Pt(0, (resized.Bounds().Dy()-headH)/2)
		}

		draw.Draw(rootImg, head, resized, newSp, draw.Src)
	}

	ximg := xgraphics.NewConvert(xw.conn, rootImg)
	defer ximg.Destroy()

	if err := ximg.XSurfaceSet(ximg.X.RootWin()); err != nil {
		return err
	}

	ximg.XDraw()
	ximg.XPaint(ximg.X.RootWin())

	return nil
}
