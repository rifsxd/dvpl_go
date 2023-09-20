//go:generate goversioninfo

package main

import (
	"errors"
	"flag"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/pierrec/lz4/v4"
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
const Name = "DVPLGO CLI CONVERTER"
const Version = "3.6.0"
const Repo = "https://github.com/RifsxD/dvpl-go"
const Web = "https://rxd-mods.xyz"
const Build = "20/09/2023"
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
	fmt.Println(`dvplgo [-mode] [-keep-originals] [-path]

    • mode can be one of the following:

        compress: compresses files into dvpl.
        decompress: decompresses dvpl files into standard files.
        help: show this help message.

	• flags can be one of the following:

    	-keep-originals flag keeps the original files after compression/decompression.
		-path specifies the directory/files path to process. Default is the current directory.

	• usage can be one of the following examples:

		$ dvplgo -mode help

		$ dvplgo -mode decompress -path /path/to/decompress/compress
		
		$ dvplgo -mode compress -path /path/to/decompress/compress
		
		$ dvplgo -mode decompress -keep-originals -path /path/to/decompress/compress
		
		$ dvplgo -mode compress -keep-originals -path /path/to/decompress/compress
		
		$ dvplgo -mode decompress -path /path/to/decompress/compress.yaml.dvpl
		
		$ dvplgo -mode compress -path /path/to/decompress/compress.yaml
		
		$ dvplgo -mode decompress -keep-originals -path /path/to/decompress/compress.yaml.dvpl
		
		$ dvplgo -mode dcompress -keep-originals -path /path/to/decompress/compress.yaml
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
				processedBlock, err = CompressDVPL(fileData)
				newName = directoryOrFile + ".dvpl"
			} else {
				processedBlock, err = DecompressDVPL(fileData)
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

// CompressDVPL compresses a buffer and returns the processed DVPL file buffer.
func CompressDVPL(buffer []byte) ([]byte, error) {
	compressedBlockSize := lz4.CompressBlockBound(len(buffer))
	compressedBlock := make([]byte, compressedBlockSize)

	n, err := lz4.CompressBlock(buffer, compressedBlock, nil)
	if err != nil {
		return nil, err
	}

	compressedBlock = compressedBlock[:n]
	footerBuffer := createDVPLFooter(uint32(len(buffer)), uint32(n), crc32.ChecksumIEEE(compressedBlock), 2)
	return append(compressedBlock, footerBuffer...), nil
}

// DecompressDVPL decompresses a DVPL buffer and returns the uncompressed file buffer.
func DecompressDVPL(buffer []byte) ([]byte, error) {
	footerData, err := readDVPLFooter(buffer)
	if err != nil {
		return nil, err
	}

	targetBlock := buffer[:len(buffer)-20]

	if uint32(len(targetBlock)) != footerData.CompressedSize {
		return nil, errors.New("DVPLSizeMismatch")
	}

	if crc32.ChecksumIEEE(targetBlock) != footerData.CRC32 {
		return nil, errors.New("DVPLCRC32Mismatch")
	}

	if footerData.Type == 0 {
		if !(footerData.OriginalSize == footerData.CompressedSize && footerData.Type == 0) {
			return nil, errors.New("DVPLTypeSizeMismatch")
		}
		return targetBlock, nil
	} else if footerData.Type == 1 || footerData.Type == 2 {
		deDVPLBlock := make([]byte, footerData.OriginalSize)
		n, err := lz4.UncompressBlock(targetBlock, deDVPLBlock)
		if err != nil {
			return nil, err
		}

		if uint32(n) != footerData.OriginalSize {
			return nil, errors.New("DVPLDecodeSizeMismatch")
		}

		return deDVPLBlock, nil
	}

	return nil, errors.New("UNKNOWN DVPL FORMAT")
}

// createDVPLFooter creates a DVPL footer from the provided data.
func createDVPLFooter(inputSize, compressedSize, crc32, typeVal uint32) []byte {
	result := make([]byte, 20)
	writeLittleEndianUint32(result, inputSize, 0)
	writeLittleEndianUint32(result, compressedSize, 4)
	writeLittleEndianUint32(result, crc32, 8)
	writeLittleEndianUint32(result, typeVal, 12)
	copy(result[16:], "DVPL")
	return result
}

func readLittleEndianUint32(b []byte, offset int) uint32 {
	return uint32(b[offset]) | uint32(b[offset+1])<<8 | uint32(b[offset+2])<<16 | uint32(b[offset+3])<<24
}

// readDVPLFooter reads the DVPL footer data from a DVPL buffer.
func readDVPLFooter(buffer []byte) (*DVPLFooter, error) {
	footerBuffer := buffer[len(buffer)-20:]
	if string(footerBuffer[16:]) != "DVPL" || len(footerBuffer) != 20 {
		return nil, errors.New(RedColor + "InvalidDVPLFooter" + ResetColor)
	}

	footerData := &DVPLFooter{}
	footerData.OriginalSize = readLittleEndianUint32(footerBuffer, 0)
	footerData.CompressedSize = readLittleEndianUint32(footerBuffer, 4)
	footerData.CRC32 = readLittleEndianUint32(footerBuffer, 8)
	footerData.Type = readLittleEndianUint32(footerBuffer, 12)
	return footerData, nil
}

func writeLittleEndianUint32(b []byte, v uint32, offset int) {
	b[offset+0] = byte(v)
	b[offset+1] = byte(v >> 8)
	b[offset+2] = byte(v >> 16)
	b[offset+3] = byte(v >> 24)
}
