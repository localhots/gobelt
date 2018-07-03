package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"reflect"
	"strings"
)

type genericType struct {
	name     string
	testVals []string
}

var types = []genericType{
	// Int
	{name: reflect.Int.String(), testVals: []string{"1", "2", "3", "4", "5"}},
	{name: reflect.Int8.String(), testVals: []string{"1", "2", "3", "4", "5"}},
	{name: reflect.Int16.String(), testVals: []string{"1", "2", "3", "4", "5"}},
	{name: reflect.Int32.String(), testVals: []string{"1", "2", "3", "4", "5"}},
	{name: reflect.Int64.String(), testVals: []string{"1", "2", "3", "4", "5"}},
	// Uint
	{name: reflect.Uint.String(), testVals: []string{"1", "2", "3", "4", "5"}},
	{name: reflect.Uint8.String(), testVals: []string{"1", "2", "3", "4", "5"}},
	{name: reflect.Uint16.String(), testVals: []string{"1", "2", "3", "4", "5"}},
	{name: reflect.Uint32.String(), testVals: []string{"1", "2", "3", "4", "5"}},
	{name: reflect.Uint64.String(), testVals: []string{"1", "2", "3", "4", "5"}},
	// String
	{name: reflect.String.String(), testVals: []string{`"1"`, `"2"`, `"3"`, `"4"`, `"5"`}},
}

func main() {
	tplDir := flag.String("tpl", "", "Path to template directory")
	destDir := flag.String("dest", "", "Path to destination directory")
	flag.Parse()

	if *tplDir == "" {
		log.Println("Template directory is not specified")
		flag.Usage()
		os.Exit(1)
	}
	if *destDir == "" {
		log.Println("Destination directory is not specified")
		flag.Usage()
		os.Exit(1)
	}

	implSource, err := ioutil.ReadFile(path.Join(*tplDir, "set.go"))
	if err != nil {
		log.Fatalf("Failed to read source file at %s: %v", path.Join(*tplDir, "set.go"), err)
	}

	testSource, err := ioutil.ReadFile(path.Join(*tplDir, "set_test.go"))
	if err != nil {
		log.Fatalf("Failed to read test file at %s: %v", path.Join(*tplDir, "set_test.go"), err)
	}

	for _, typ := range types {
		log.Printf("Generating package for type %s\n", typ.name)
		err := generate(*destDir, implSource, testSource, typ)
		if err != nil {
			log.Fatalf("Failed to generate a package for type %s: %v", typ.name, err)
		}
	}
	err = gofmt(*destDir)
	if err != nil {
		log.Fatalf("Formatting failed: %v", err)
	}

	log.Println("Set packages were successfully generated")
}

func generate(destDir string, implSource, testSource []byte, typ genericType) error {
	pkgDir := path.Join(destDir, fmt.Sprintf("set%s", typ.name))
	err := os.RemoveAll(pkgDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = os.Mkdir(pkgDir, 0777)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(pkgDir, "set.go"), renderBytes(implSource, typ), 0755)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(pkgDir, "set_test.go"), renderBytes(testSource, typ), 0755)
	return err
}

func renderBytes(src []byte, typ genericType) []byte {
	return []byte(render(string(src), typ))
}

func render(src string, typ genericType) string {
	const genericTypeName = "TypeName"
	const doNotEditBanner = `/*******************************************************************************
THIS FILE WAS AUTOMATICALLY GENERATED. DO NOT EDIT!
*******************************************************************************/
`
	// Replace test constants
	src = strings.Replace(src, "One   TypeName = 1", "One   TypeName = "+typ.testVals[0], 1)
	src = strings.Replace(src, "Two   TypeName = 2", "Two   TypeName = "+typ.testVals[1], 1)
	src = strings.Replace(src, "Three TypeName = 3", "Three TypeName = "+typ.testVals[2], 1)
	src = strings.Replace(src, "Four  TypeName = 4", "Four  TypeName = "+typ.testVals[3], 1)
	src = strings.Replace(src, "Five  TypeName = 5", "Five  TypeName = "+typ.testVals[4], 1)
	// Replace the type name
	src = strings.Replace(src, genericTypeName, typ.name, -1)
	// Replace package name
	src = strings.Replace(src, "package impl", "package set"+typ.name, 1)
	return doNotEditBanner + src
}

func gofmt(dir string) error {
	out, err := exec.Command("gofmt", "-w", "-l", dir).CombinedOutput()
	if err != nil {
		log.Println("gofmt returned:", string(out))
	}
	return err
}
