package main

import (
	"context"
	"flag"
	"image/png"
	"log"
	"os"

	"github.com/sfomuseum/go-activitypub/icon"
)

func main() {

	var label string
	var trim_to int
	var font_size float64

	var outfile string

	flag.StringVar(&label, "label", "", "...")
	flag.IntVar(&trim_to, "trim-to", 0, "...")
	flag.StringVar(&outfile, "out", "icon.png", "...")
	flag.Float64Var(&font_size, "font-size", 48.0, "...")

	flag.Parse()

	ctx := context.Background()

	opts := &icon.GenerateIconOptions{
		Label:    label,
		TrimTo:   trim_to,
		FontSize: font_size,
	}

	im, err := icon.GenerateIcon(ctx, opts)

	if err != nil {
		log.Fatalf("Failed to generate icon")
	}

	wr, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalf("Failed to open '%s' for writing, %v", outfile, err)
	}

	err = png.Encode(wr, im)

	if err != nil {
		log.Fatalf("Failed to encode '%s', %v", outfile, err)
	}

	err = wr.Close()

	if err != nil {
		log.Fatalf("Failed to close '%s' after writing, %v", outfile, err)
	}

}
