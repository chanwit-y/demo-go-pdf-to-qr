package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/multi/qrcode"
)

func convertPDFToImage(pdfPath string, imagePath string) error {
	// Get the absolute paths of the PDF and image files
	absPDFPath, err := filepath.Abs(pdfPath)
	if err != nil {
		return err
	}
	absImagePath, err := filepath.Abs(imagePath)
	if err != nil {
		return err
	}

	// Define the Ghostscript command to convert PDF to image
	cmd := exec.Command("gs", "-dQUIET", "-dSAFER", "-dBATCH", "-dNOPAUSE", "-sDEVICE=jpeg", "-r144", "-sOutputFile="+absImagePath, absPDFPath)

	// Execute the Ghostscript command
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil

}

func readImage(imagePath string) ([]byte, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func scan(b []byte) (string, string) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		msg := fmt.Sprintf("failed to read image: %v", err)
		return "", msg
	}

	source := gozxing.NewLuminanceSourceFromImage(img)
	bin := gozxing.NewHybridBinarizer(source)
	bbm, err := gozxing.NewBinaryBitmap(bin)

	if err != nil {
		msg := fmt.Sprintf("error during processing: %v", err)
		return "", msg
	}

	qrReader := qrcode.NewQRCodeMultiReader()
	result, err := qrReader.DecodeMultiple(bbm, nil)
	if err != nil {
		msg := fmt.Sprintf("unable to decode QRCode: %v", err)
		return "", msg
	}
	strRes := []string{}
	for _, element := range result {
		strRes = append(strRes, element.String())
	}

	res := strings.Join(strRes, "\n")
	return res, ""
}

func main() {
	// pdfPath := "example/APV-122060130-1855.pdf"
	pdfPath := "example/PV-923050005-3857.pdf"
	imagePath := "image/APV-122060130-1855.jpg"

	err := convertPDFToImage(pdfPath, imagePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("PDF converted to image successfully.")
	b, _ := readImage(imagePath)
	r, _ := scan(b)

	fmt.Printf("QR: %v \n", r)
}
