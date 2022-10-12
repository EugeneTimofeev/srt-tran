package main

import (
	"bufio"
	"fmt"
	"github.com/Conight/go-googletrans"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type linesForTranslate struct {
	numInRange int64
	stringsEn  string
	stringsRu  string
}

type fileForTranslate struct {
	fileName  string
	fileLines []linesForTranslate
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("No args")
		os.Exit(0)
	}
	for _, cat := range args {		
		var filesForTranslate []fileForTranslate
		getSRT(&filesForTranslate, cat)
		for _, v := range filesForTranslate {
			doMake(&v)
		}
	}
}

func doMake(fft *fileForTranslate) {
	t := translator.New()
	fft.makeFileForTranslate()
	for j := 0; j < len(fft.fileLines); j++ {
		result, err := t.Translate(fft.fileLines[j].stringsEn, "en", "ru")
		if err != nil {
			log.Fatal(err)
		} else {
			fft.fileLines[j].stringsRu = result.Text
		}
	}
	fft.saveRuFiles()
}

func (fft *fileForTranslate) makeFileForTranslate() {
	var i int64 = 0
	var str string
	f, err := os.Open(fft.fileName)
	if err != nil {
		log.Fatal(err)
	}
	fft.fileLines = append(fft.fileLines, linesForTranslate{i, "", ""})
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		str = scanner.Text()
		if len(fft.fileLines[i].stringsEn)+len(str) <= 4000 {
			fft.fileLines[i].stringsEn = fft.fileLines[i].stringsEn + "\n" + str
		} else {
			fft.fileLines[i].stringsEn = fft.fileLines[i].stringsEn + "\n"
			i++
			fft.fileLines = append(fft.fileLines, linesForTranslate{i, str, ""})
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func (fft *fileForTranslate) saveRuFiles() {
	ruFile := fft.fileName[0:len(fft.fileName)-4] + "_rus.srt"
	_, err := os.Create(ruFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	f, err := os.OpenFile(ruFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	for i := 0; i < len(fft.fileLines); i++ {
		_, err = fmt.Fprintln(f, fft.fileLines[i].stringsRu)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Done", ruFile)
}

func getSRT(fft *[]fileForTranslate, cat string) {
	err := filepath.Walk(cat,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if (!info.IsDir()) && (filepath.Ext(path) == ".srt") && (!strings.HasSuffix(path, "_rus.srt")) {
				*fft = append(*fft, fileForTranslate{path, *new([]linesForTranslate)})
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
