package main

import (
	"fmt"
	"github.com/oliamb/cutter"
	"github.com/signintech/gopdf"
	"image"
	"image/png"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fileName := "/tmp/test.png"
	border := 20

	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	fconf, err := os.Open(fileName)
	defer fconf.Close()
	if err != nil {
		log.Fatal(err)
	}

	cfg, _, err := image.DecodeConfig(fconf)
	if err != nil {
		log.Fatal(err)
	}

	h := cfg.Width * 297 / 210

	chunks := cfg.Height / h
	if cfg.Height%h > 0 {
		chunks++
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	a4Size := gopdf.PageSizeA4

	for i := 0; i < chunks; i++ {
		tmpFile := fmt.Sprintf("/tmp/cutter/cutted%d.png", i)
		cImg, err := cutter.Crop(img, cutter.Config{
			Anchor: struct{ X, Y int }{X: 0, Y: h * i},
			Width:  cfg.Width,
			Height: h,
		})

		if err != nil {
			log.Fatal(err)
		}

		fo, err := os.Create(tmpFile)
		if err != nil {
			log.Fatal("Cannot create output file:", err)
		}

		err = png.Encode(fo, cImg)
		fo.Close()

		pdf.AddPage()

		sizeT := *a4Size
		size := &sizeT

		if i == chunks-1 {
			size.H = size.W * float64(cfg.Height%h) / float64(cfg.Width)
		}

		size.H -= float64(border)
		size.W -= float64(border)

		err = pdf.Image(tmpFile, float64(border/2), float64(border/2), size)
		if err != nil {
			log.Fatal("Cannot create output file:", err)
		}

		go func(name string) {
			log.Println("Removing ", name)
			err = os.Remove(name)
			if err != nil {
				log.Fatal(err)
			}
		}(tmpFile)
	}

	pdf.WritePdf("/tmp/cutter/cutted.pdf")

}
