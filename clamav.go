package main

import (
	"bytes"
	"fmt"
	"github.com/crazytyper/clamav"
	"io"
)

func initClamAV() (*clamav.Engine, error) {
	var err error
	if err = clamav.Init(clamav.InitDefault); err != nil {
		return nil, err
	}
	engine := clamav.New()
	_, err = engine.Load(clamav.DBDir(), clamav.DbStdopt)
	if err != nil {
		return nil, err
	}
	if err = engine.Compile(); err != nil {
		return nil, err
	}
	return engine, err
}

type scanFunc func(io.Reader) (string, error)

func CreateInMemoryFileScanner() (scanFunc, error) {
	engine, err := initClamAV()
	if err != nil {
		return nil, err
	}
	opts := clamav.ScanOptions{
		clamav.ScanGeneralAllmatches | clamav.ScanGeneralHeuristics,
		clamav.ScanParsePdf,
		0,
		0,
		0,
	}
	scanner := func(data io.Reader) (string, error) {
		fmt.Println("scanning file")
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(data)
		if err != nil {
			return "", err
		}
		fmap := clamav.OpenMemory(buf.Bytes())
		defer clamav.CloseMemory(fmap)
		virus, _, err := engine.ScanMapCb(fmap, "upload", &opts, "scan")
		// if virus != "" {
		// 	fmt.Printf("Found virus %+v\n", virus)
		// } else {
		// 	fmt.Println("No virus")
		// }
		return virus, err
	}
	return scanner, nil
}
