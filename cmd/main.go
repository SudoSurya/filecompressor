package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
)

type FileInfo struct {
	name     string
	size     int64
	contents string
	filetype string
}

func main() {
	args := os.Args

	FILENAME := args[1]

	fmt.Println(filepath.Ext(FILENAME))

	switch filepath.Ext(FILENAME) {
	case ".jpeg", ".jpg":
        err := compressJPEG(FILENAME)
        if err != nil {
            fmt.Println(err)
        }
	}
}

func compressJPEG(FILENAME string) error {
	file, err := os.Open(FILENAME)
	if err != nil {
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()

	if err != nil {
		return err
	}
	originalFILE := FileInfo{
		name:     FILENAME,
		size:     fileStat.Size(),
		filetype: filepath.Ext(FILENAME),
	}
	fmt.Println(originalFILE)

	img, err := jpeg.Decode(file)

	if err != nil {
		return err
	}

	newFILE, err := os.Create("compressed_" + FILENAME)

	if err != nil {
		return err
	}

	defer newFILE.Close()

	err = jpeg.Encode(newFILE, img, &jpeg.Options{Quality: 70})

	if err != nil {
		return err
	}

	newFILEStat, err := newFILE.Stat()

	if err != nil {
		return err
	}

	newFILEINFO := FileInfo{
		name:     "compressed_" + FILENAME,
		size:     newFILEStat.Size(),
		filetype: filepath.Ext(FILENAME),
	}

	fmt.Println(newFILEINFO)

	return nil
}

/* func readFileContents(FILENAME string) (FileInfo, error) {
    fileType := filepath.Ext(FILENAME)
    file, err := os.Open(FILENAME)

    if err != nil {
        return FileInfo{}, err
    }

    defer file.Close()

    fileStat, err := file.Stat()

    if err != nil {
        return FileInfo{}, err
    }

    fileContents := make([]byte, fileStat.Size())

    _, err = file.Read(fileContents)

    if err != nil {
        return FileInfo{}, err
    }

    return FileInfo{FILENAME, fileStat.Size(), string(fileContents),fileType}, nil

} */
