package main

import (
	"math/rand"
	"os"
	"time"
	lession_5 "tinyrender-golang/lesson/lession-5"
	model "tinyrender-golang/model"
	"tinyrender-golang/tga"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {

	f, err := os.Open("./obj/african_head/african_head.obj")
	if err != nil {
		panic(err)
	}
	obj, err := model.NewReader(f).Read()
	if err != nil {
		panic(err)
	}
	textureFile, err := os.Open("./obj/african_head/african_head_diffuse.tga")
	if err != nil {
		panic(err)
	}
	texture, err := tga.DecodeToTga(textureFile)
	if err != nil {
		panic(err)
	}
	texture.FlipVertical()

	fb := tga.CreateTga(800, 800)

	lession_5.DrawWithCamera(fb, obj, texture)

	fb.FlipVertical()
	err = fb.SaveToFile("./pic/lesson-5-1.tga")
	if err != nil {
		panic(err)
	}
}
