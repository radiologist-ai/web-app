package app

import (
	"bytes"
	"context"
	"fmt"
	"github.com/radiologist-ai/web-app/internal/domain"
	"os"
)

func DoShit(ctx context.Context, mc domain.S3Adapter) {
	fname := "video.mp4"
	file, err := os.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	chunkSize := 5 << 20
	size := len(file)
	l := 0
	r := chunkSize
	i := 1
	fmt.Println("read file '"+fname+"' size = ", size)
	fmt.Println("chunks num", (size+chunkSize-1)/chunkSize)
	savedObjs := make([]string, 0)
	for {
		if r > size {
			r = size
		}
		chunk := file[l:r:r]

		data := bytes.NewReader(chunk)
		objName, err := mc.Save(ctx, fmt.Sprintf("video/chunk%d", i), "application/octet-stream", data, int64(len(chunk)))
		if err != nil {
			panic(err)
		}
		savedObjs = append(savedObjs, objName)
		if r == size {
			break
		}
		l += chunkSize
		r += chunkSize
		i++
	}
	fmt.Println("saved all chunks", savedObjs)
	fmt.Println("start composing objects")
	result, err := mc.Compose(ctx, fname, savedObjs...)
	if err != nil {
		panic(err)
	}
	fmt.Println("finished composing object + '" + result + "'")
}
