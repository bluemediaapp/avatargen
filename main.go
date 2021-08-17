package main

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/gofiber/fiber/v2"
)

var (
	app        = fiber.New()
	config     *Config
	assetTypes = []string{"beak", "body", "hat", "feet", "wing"}
	assets     = make(map[string][]image.Image)
)

func main() {
	config = &Config{
		port: os.Getenv("PORT"),
	}
	discoverAssets()
	app.Post("/", func(ctx *fiber.Ctx) error {
		user_id := ctx.Get("user_id")
		assetIds := md5.Sum([]byte(user_id))

		duck := gg.NewContext(499, 600)

		for i, assetType := range assetTypes {
			assetId := uint16(assetIds[i])
			asset, err := getAsset(assetId, assetType)
			if err != nil {
				continue
			}
			duck.DrawImage(asset, 0, 0)
		}

		// Exporting
		out := bytes.NewBuffer(nil)
		duck.EncodePNG(out)
		ctx.Set("Content-Type", "image/png")
		ctx.Write(out.Bytes())

		return nil
	})
	log.Print(app.Listen(config.port))
}

func getAsset(assetId uint16, assetType string) (image.Image, error) {
	assetCounts := uint16(len(assets[assetType]))
	if assetCounts <= 1 {
		return nil, errors.New("Not enough assets in category %s")
	}
	actualAssetId := assetId % assetCounts // Wrap around (AID 7, AC 5 = 0 1 2 3 4 0 (1)) where one becomes the output
	return assets[assetType][actualAssetId], nil
}

func discoverAssets() {
	for _, assetType := range assetTypes {
		pendingAssets := make([]image.Image, 0)
		files, err := ioutil.ReadDir(fmt.Sprintf("assets/%s", assetType))
		if err != nil {
			log.Printf("Folder read error")
			log.Print(err)
		}
		for _, file := range files {
			if !strings.HasSuffix(file.Name(), ".png") {
				continue
			}
			image, err := gg.LoadImage(fmt.Sprintf("assets/%s/%s", assetType, file.Name()))
			if err != nil {
				log.Print(err)
				continue
			}
			pendingAssets = append(pendingAssets, image)
		}
		log.Printf("Loading %d for category %s", len(pendingAssets), assetType)
		assets[assetType] = pendingAssets
	}
}
