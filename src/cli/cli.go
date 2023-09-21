//go:generate goversioninfo

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/rifsxd/dvpl_go/dvpl_logic"
)

const (
	dvplFooterSize = 20
	dvplTypeNone   = 0
	dvplTypeLZ4    = 2
	dvplExtension  = ".dvpl"
)

// ANSI escape codes for text coloring
const (
	RedColor    = "\033[31m"
	GreenColor  = "\033[32m"
	YellowColor = "\033[33m"
	ResetColor  = "\033[0m"
)

// Config represents the configuration for the program.
type Config struct {
	Mode          string
	KeepOriginals bool
	Path          string // New field to specify the directory path.
}

// DVPLFooter represents the DVPL file footer data.
type DVPLFooter struct {
	OriginalSize   uint32
	CompressedSize uint32
	CRC32          uint32
	Type           uint32
}

// Info variables
const Dev = "RifsxD"
const Name = "DVPL_GO CLI CONVERTER"
const Version = "4.1.0"
const Repo = "https://github.com/RifsxD/dvpl-go"
const Web = "https://rxd-mods.xyz"
const Build = "21/09/2023"
const Info = "A CLI Tool Coded In JavaScript To Convert WoTB ( Dava ) SmartDLC DVPL File Based On LZ4_HC Compression."

func main() {

	cyan := color.New(color.FgCyan)

	fmt.Println()
	cyan.Println("• Name:", Name)
	cyan.Println("• Version:", Version)
	cyan.Println("• Build:", Build)
	cyan.Println("• Dev:", Dev)
	cyan.Println("• Repo:", Repo)
	cyan.Println("• Web:", Web)
	cyan.Println("• Info:", Info)
	fmt.Println()

	config, err := parseCommandLineArgs()
	if err != nil {
		log.Printf("%sError%s parsing command-line arguments: %v", RedColor, ResetColor, err)
		return
	}

	switch config.Mode {
	case "compress", "decompress":
		err := processFiles(config.Path, config)
		if err != nil {
			log.Printf("%s%s FAILED%s: %v", RedColor, strings.ToUpper(config.Mode), ResetColor, err)
		} else {
			log.Printf("%s%s FINISHED%s.", GreenColor, strings.ToUpper(config.Mode), ResetColor)
		}
	case "help":
		printHelpMessage()
	default:
		log.Fatalf("%sIncorrect mode selected. Use '-help' for information.%s", RedColor, ResetColor)
	}
}

func parseCommandLineArgs() (*Config, error) {
	config := &Config{}
	flag.StringVar(&config.Mode, "mode", "", "Mode can be 'compress' / 'decompress' / 'help' (for an extended help guide).")
	flag.BoolVar(&config.KeepOriginals, "keep-originals", false, "Keep original files after compression/decompression.")
	flag.StringVar(&config.Path, "path", ".", "directory/files path to process. Default is the current directory.")
	flag.Parse()

	if config.Mode == "" {
		return nil, errors.New("No mode selected. Use '-help' for usage information.")
	}

	return config, nil
}

func printHelpMessage() {
	fmt.Println(`dvpl_go [-mode] [-keep-originals] [-path]

    • mode can be one of the following:

        compress: compresses files into dvpl.
        decompress: decompresses dvpl files into standard files.
        help: show this help message.

	• flags can be one of the following:

    	-keep-originals flag keeps the original files after compression/decompression.
		-path specifies the directory/files path to process. Default is the current directory.

	• usage can be one of the following examples:

		$ dvpl_go -mode help

		$ dvpl_go -mode decompress -path /path/to/decompress/compress
		
		$ dvpl_go -mode compress -path /path/to/decompress/compress
		
		$ dvpl_go -mode decompress -keep-originals -path /path/to/decompress/compress
		
		$ dvpl_go -mode compress -keep-originals -path /path/to/decompress/compress
		
		$ dvpl_go -mode decompress -path /path/to/decompress/compress.yaml.dvpl
		
		$ dvpl_go -mode compress -path /path/to/decompress/compress.yaml
		
		$ dvpl_go -mode decompress -keep-originals -path /path/to/decompress/compress.yaml.dvpl
		
		$ dvpl_go -mode dcompress -keep-originals -path /path/to/decompress/compress.yaml
	`)
}

func processFiles(directoryOrFile string, config *Config) error {
	info, err := os.Stat(directoryOrFile)
	if err != nil {
		return err
	}

	if info.IsDir() {
		dirList, err := os.ReadDir(directoryOrFile)
		if err != nil {
			return err
		}

		for _, dirItem := range dirList {
			err := processFiles(filepath.Join(directoryOrFile, dirItem.Name()), config)
			if err != nil {
				fmt.Printf("Error processing directory %s: %v\n", dirItem.Name(), err)
			}
		}
	} else {
		isDecompression := config.Mode == "decompress" && strings.HasSuffix(directoryOrFile, ".dvpl")
		isCompression := config.Mode == "compress" && !strings.HasSuffix(directoryOrFile, ".dvpl")

		if isDecompression || isCompression {
			filePath := directoryOrFile
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("%sError%s reading file %s: %v\n", RedColor, ResetColor, directoryOrFile, err)
				return err
			}

			var processedBlock []byte
			newName := ""

			if isCompression {
				processedBlock, err = dvpl_logic.CompressDVPL(fileData)
				newName = directoryOrFile + ".dvpl"
			} else {
				processedBlock, err = dvpl_logic.DecompressDVPL(fileData)
				newName = strings.TrimSuffix(directoryOrFile, ".dvpl")
			}

			if err != nil {
				fmt.Printf("File %s failed to convert due to %v\n", directoryOrFile, err)
				return err
			}

			err = os.WriteFile(newName, processedBlock, 0644)
			if err != nil {
				fmt.Printf("%sError%s writing file %s: %v\n", RedColor, ResetColor, newName, err)
				return err
			}

			fmt.Printf("File %s has been successfully %s into %s%s%s\n", filePath, getAction(config.Mode), GreenColor, newName, ResetColor)

			if !config.KeepOriginals {
				err := os.Remove(filePath)
				if err != nil {
					fmt.Printf("%sError%s deleting file %s: %v\n", RedColor, ResetColor, filePath, err)
				}
			}
		} else {
			fmt.Printf("%sIgnoring%s file %s\n", YellowColor, ResetColor, directoryOrFile)
		}
	}

	return nil
}

func getAction(mode string) string {
	if mode == "compress" {
		return GreenColor + "compressed" + ResetColor
	}
	return GreenColor + "decompressed" + ResetColor
}
