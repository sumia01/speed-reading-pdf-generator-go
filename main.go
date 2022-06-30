package main

import (
	_ "embed"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/signintech/gopdf"
)

//go:embed font/roboto-regular.ttf
var font []byte

var pageSize = &gopdf.Rect{H: 420, W: 595}

const (
	fontName = "roboto-regular"
	sep      = string(os.PathSeparator)
	inDir    = "." + sep + "In"
	outDir   = "." + sep + "Out"
	pdfExt   = ".pdf"
	txtExt   = ".txt"
	fontSize = 30
)

func main() {
	_, err := os.Stat(inDir)
	if os.IsNotExist(err) {
		log.Fatalf("%v folder is not exists, please create it and put your txt files there.", inDir)
	}

	err = createFolderOrFailifNotExists(outDir)
	if err != nil {
		log.Fatalf("%v is not exits and can't create it: %v", outDir, err)
	}

	fileinfos, err := ioutil.ReadDir(inDir)
	if err != nil {
		log.Fatal("can't read files from \"In\" folder, make sure you created it.")
	}

	for _, fileinfo := range fileinfos {
		if fileinfo.IsDir() {
			continue
		}

		if strings.HasSuffix(strings.ToLower(fileinfo.Name()), txtExt) {
			file, err := ioutil.ReadFile(inDir + sep + fileinfo.Name())
			if err != nil {
				log.Printf("can't read %v from \"In\" folder, skip generating pdf from it,  error: %v\n", fileinfo.Name(), err)
				continue
			}

			words := strings.Split(string(file), " ")
			err = createPdfFromText(
				words,
				strings.TrimSuffix(
					strings.ToLower(fileinfo.Name()),
					txtExt,
				)+pdfExt,
			)
			if err != nil {
				log.Printf("can't generate pdf from %v: error: %v\n", fileinfo.Name(), err)
			}
		}
	}
}

func createPdfFromText(words []string, filename string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *pageSize})
	err := pdf.AddTTFFontData(fontName, font)
	if err != nil {
		return err
	}

	err = pdf.SetFont(fontName, "", fontSize)
	if err != nil {
		return err
	}

	for _, word := range words {
		size := pageSize
		textwidth, err := pdf.MeasureTextWidth(word)
		if err != nil {
			continue
		}

		pdf.SetMargins(size.W/2-(textwidth/2), size.H/2-(fontSize/2), 0, 0)

		if len(word) == 0 {
			continue
		}

		pdf.AddPage()

		err = pdf.Cell(nil, word)
		if err != nil {
			continue
		}

	}

	return pdf.WritePdf(outDir + sep + filename)
}

func createFolderOrFailifNotExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, os.ModePerm)
		if errDir != nil {
			return errDir
		}
	}
	return nil
}
