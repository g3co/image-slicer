package cmd

import (
	"fmt"
	"github.com/oliamb/cutter"
	"github.com/signintech/gopdf"
	"github.com/spf13/cobra"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

func Execute() {

	border := 0
	outputFile := ""

	RootCmd := &cobra.Command{
		Use:              "image-slicer [file name] -b=[border size] -o=[output file]",
		Short:            "Image slicer",
		Long:             `This application slices image to the pdf bundle`,
		TraverseChildren: true,
		Args:             cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			fileName := args[0]

			fileParts := strings.Split(fileName, ".")
			fileType := fileParts[len(fileParts)-1]

			if outputFile == "" {
				fileParts[len(fileParts)-1] = "pdf"
				outputFile = strings.Join(fileParts, ".")
			}

			imgSourceFile, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			defer imgSourceFile.Close()

			img, _, err := image.Decode(imgSourceFile)
			if err != nil {
				log.Fatal(err)
			}

			cfgSourceFile, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			defer cfgSourceFile.Close()

			cfg, _, err := image.DecodeConfig(cfgSourceFile)
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
				tmpFile := fmt.Sprintf("%s_tmp_%d", fileName, i)
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
					log.Fatal(err)
				}

				defer fo.Close()
				defer os.Remove(tmpFile)

				switch fileType {
				case "jpeg", "jpg":
					err = jpeg.Encode(fo, cImg, &jpeg.Options{100})
				case "png":
					err = png.Encode(fo, cImg)
				default:
					log.Fatal("unsupported image format:", fileType)
				}

				if err != nil {
					log.Fatal(err)
				}

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
					log.Fatal(err)
				}
			}

			pdf.WritePdf(outputFile)
		},
	}

	RootCmd.Flags().IntVarP(&border, "border", "b", 0, "border size")
	RootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
