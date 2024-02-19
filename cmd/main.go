package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"

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

func main() {
	args := os.Args

	FILENAME := args[1]

	switch filepath.Ext(FILENAME) {
	case ".jpeg", ".jpg":
		originalFile, compressFile, err := compressJPEG(FILENAME)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color("#1fb55b")).Render("Original File Info"))
		renderFileInfo(originalFile)
		fmt.Println(lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color("#f70202")).Render("Compressed File Info"))
        renderFileInfo(compressFile)
	case ".pdf":
		fmt.Println("PDF")
	}
}

func renderFileInfo(fileInfo FileInfo) {
	fmt.Print("\n")
	fmt.Println(titleStyle.Render("Name: ") + itemStyle.Render(fileInfo.name))
	fmt.Println(titleStyle.Render("Size: ") + itemStyle.Render(fmt.Sprintf("%d bytes", fileInfo.size)))
	fmt.Println(titleStyle.Render("Type: ") + itemStyle.Render(fileInfo.filetype))
	fmt.Print("\n")
}

func compressJPEG(FILENAME string) (FileInfo, FileInfo, error) {
	file, err := os.Open(FILENAME)
	if err != nil {
		return FileInfo{}, FileInfo{}, err
	}
	defer file.Close()

	fileStat, err := file.Stat()

	if err != nil {
		return FileInfo{}, FileInfo{}, err
	}
	originalFILE := FileInfo{
		name:     FILENAME,
		size:     fileStat.Size(),
		filetype: filepath.Ext(FILENAME),
	}

	img, err := jpeg.Decode(file)

	if err != nil {
		return FileInfo{}, FileInfo{}, err
	}

	compressedFILE, err := os.Create("compressed_" + FILENAME)

	if err != nil {
		return FileInfo{}, FileInfo{}, err
	}

	defer compressedFILE.Close()

	err = jpeg.Encode(compressedFILE, img, &jpeg.Options{Quality: 70})

	if err != nil {
		return FileInfo{}, FileInfo{}, err
	}

	compressedFILEStat, err := compressedFILE.Stat()

	if err != nil {
		return FileInfo{}, FileInfo{}, err
	}

	compressedFILEINFO := FileInfo{
		name:     "compressed_" + FILENAME,
		size:     compressedFILEStat.Size(),
		filetype: filepath.Ext(FILENAME),
	}


	return originalFILE, compressedFILEINFO, nil
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
