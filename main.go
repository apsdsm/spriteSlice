package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) < 4 {
		fmt.Println("usage: spriteSplice [imagename] [tile width] [tile height]")
		os.Exit(0)
	}

	// get args
	pathToImg := os.Args[1]
	tileW, _ := strconv.Atoi(os.Args[2])
	tileH, _ := strconv.Atoi(os.Args[3])

	// strings used to create new files
	dirName, fileName := path.Split(pathToImg)
	extName := path.Ext(fileName)
	resName := strings.TrimSuffix(fileName, extName)

	// try open image
	file, err := os.Open(pathToImg)

	if err != nil {
		fmt.Println("was not able to open file: ", pathToImg)
		os.Exit(0)
	}

	defer file.Close()

	// try decode image data
	data, _, err := image.Decode(file)

	if err != nil {
		fmt.Println("was not able to decode image: ", pathToImg)
		os.Exit(0)
	}

	// try make dir to hold new images
	expPath := path.Join(dirName, resName+"_split")
	err = os.Mkdir(expPath, 0777)

	if err != nil {
		fmt.Println("was not able to make export directory: ", expPath)
		os.Exit(0)
	}

	// variables to extract image data
	bounds := data.Bounds()
	tilesAcross := (bounds.Max.X - bounds.Min.X) / tileW
	tilesDown := (bounds.Max.Y - bounds.Min.Y) / tileH
	newImageBounds := image.Rect(0, 0, tileW, tileH)

	// for each tile
	for x := 0; x < tilesAcross; x++ {
		for y := 0; y < tilesDown; y++ {

			fmt.Println("processing tile", x, y)

			// define start/end coords for source data
			sourceXStart := x * tileW
			sourceXEnd := sourceXStart + tileW
			sourceYStart := y * tileH
			sourceYEnd := sourceYStart + tileH

			// define starting values for extraction loop
			destX := 0
			destY := 0
			sourceX := sourceXStart
			sourceY := sourceYStart

			// make new image to store data
			tile := image.NewNRGBA(newImageBounds)

			// copy source pixels to new image
			for sourceX < sourceXEnd {
				for sourceY < sourceYEnd {
					tile.Set(destX, destY, data.At(sourceX, sourceY))

					sourceY++
					destY++
				}

				destY = 0
				sourceY = sourceYStart
				sourceX++
				destX++
			}

			// create and save new image
			newImgPath := expPath + "/" + fileName + "_" + strconv.Itoa(x) + "_" + strconv.Itoa(y) + extName
			f, _ := os.Create(newImgPath)
			png.Encode(f, tile)
			f.Close()
		}
	}

	fmt.Println("done!")
}
