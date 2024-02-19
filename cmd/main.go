package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type FileInfo struct {
	name     string
	size     int64
	filetype string
}

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Margin(0, 1, 0, 1)
var itemStyle = lipgloss.NewStyle().
	Bold(true).
	Margin(0, 1, 0, 1)

var errorStyle = lipgloss.NewStyle().
	Bold(true).
	Background(lipgloss.Color("#f70202")).
	Foreground(lipgloss.Color("#ffffff")).
	Padding(1, 2)

func main() {
	args := os.Args

	FILENAME := args[1]

	switch filepath.Ext(FILENAME) {
	case ".jpeg", ".jpg":
		originalFile, err := GetFileInfo(FILENAME)
		if err != nil {
			renderError(err.Error())
		}
		renderTitle("Original File Info")
		renderFileInfo(originalFile)
		if originalFile.size < 500000 {
			renderError("File size is too small to compress")
			return
		}
		compressFile, err := compressJPEG(FILENAME)
		if err != nil {
			fmt.Println(err)
		}
		renderTitle("Compressed File Info")
		renderFileInfo(compressFile)
	case ".pdf":
		originalFile, err := GetFileInfo(FILENAME)
		if err != nil {
			renderError(err.Error())
		}
		renderTitle("Original File Info")
		renderFileInfo(originalFile)
		if originalFile.size < 500000 {
			renderError("File size is too small to compress")
			return
		}

		compressFile, err := compressPDF(FILENAME)
		if err != nil {
			renderError(err.Error())
		}
		renderTitle("Compressed File Info")
		renderFileInfo(compressFile)
	}
}

func compressPDF(FILENAME string) (FileInfo, error) {
	file, err := os.Open(FILENAME)
	if err != nil {
		return FileInfo{}, err
	}
	defer file.Close()


	FILENAME = strings.Replace(FILENAME, ".pdf", ".gz", -1)
	compressedFILE, err := os.Create("compressed_" + FILENAME)
	if err != nil {
		return FileInfo{}, err
	}
	defer compressedFILE.Close()


    gzipWriter := gzip.NewWriter(compressedFILE)
    defer gzipWriter.Close()

    _, err = io.Copy(gzipWriter, bufio.NewReader(file))
    if err != nil {
        return FileInfo{}, err
    }

	compressedFILEStat, err := compressedFILE.Stat()
	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		name:     "compressed_" + FILENAME,
		size:     compressedFILEStat.Size(),
		filetype: filepath.Ext(FILENAME),
	}, nil
}

/* compresses the provided jpeg file */
func compressJPEG(FILENAME string) (FileInfo, error) {
	file, err := os.Open(FILENAME)
	if err != nil {
		return FileInfo{}, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)

	if err != nil {
		return FileInfo{}, err
	}

	compressedFILE, err := os.Create("compressed_" + FILENAME)

	if err != nil {
		return FileInfo{}, err
	}

	defer compressedFILE.Close()

	err = jpeg.Encode(compressedFILE, img, &jpeg.Options{Quality: 70})

	if err != nil {
		return FileInfo{}, err
	}

	compressedFILEStat, err := compressedFILE.Stat()

	if err != nil {
		return FileInfo{}, err
	}

	compressedFILEINFO := FileInfo{
		name:     "compressed_" + FILENAME,
		size:     compressedFILEStat.Size(),
		filetype: filepath.Ext(FILENAME),
	}

	return compressedFILEINFO, nil
}

/* returns File Info of the provided file */
func GetFileInfo(FILENAME string) (FileInfo, error) {
	file, err := os.Open(FILENAME)
	if err != nil {
		return FileInfo{}, err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		name:     fileStat.Name(),
		size:     fileStat.Size(),
		filetype: filepath.Ext(FILENAME),
	}, nil
}

/*
utils for rendering output to the terminal
lipgloss is a library for styling terminal output
*/
func renderTitle(title string) {
	fmt.Println(lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color("#1fb55b")).Render(title))
}

func renderError(err string) {
	fmt.Println(errorStyle.Render(err))
}

func renderFileInfo(fileInfo FileInfo) {
	fmt.Print("\n")
	fmt.Println(titleStyle.Render("Name: ") + itemStyle.Render(fileInfo.name))
	fmt.Println(titleStyle.Render("Size: ") + itemStyle.Render(fmt.Sprintf("%d bytes", fileInfo.size)))
	fmt.Println(titleStyle.Render("Type: ") + itemStyle.Render(fileInfo.filetype))
	fmt.Print("\n")
}
