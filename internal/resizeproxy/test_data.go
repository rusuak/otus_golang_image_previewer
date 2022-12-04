package resizeproxy

import (
	"bufio"
	"fmt"
	"os"
)

const (
	ImageURL           = "https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/"
	SourceImgName      = "_gopher_original_1024x504.jpg"
	DestinationImgName = "gopher_256x126_resized.jpg"
)

func loadImage(imgName string) []byte {
	fileToBeUploaded := "../../img_example/" + imgName
	file, err := os.Open(fileToBeUploaded)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileInfo, _ := file.Stat()
	bytes := make([]byte, fileInfo.Size())

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	return bytes
}

func getTestImagesSupportedTypes() []string {
	return []string{"image/jpg", "image/jpeg"}
}
