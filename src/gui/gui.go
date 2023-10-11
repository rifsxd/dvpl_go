//go:generate goversioninfo -64

package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/rifsxd/dvpl_go/dvpl_logic"
)

const (
	RedColor    = "\033[31m"
	GreenColor  = "\033[32m"
	YellowColor = "\033[33m"
	ResetColor  = "\033[0m"
)

type Config struct {
	Mode          string
	KeepOriginals bool
	Path          string
}

type DVPLFooter struct {
	OriginalSize   uint32
	CompressedSize uint32
	CRC32          uint32
	Type           uint32
}

//go:embed resource/dvpl_go.png
var resources embed.FS

// Info variables
const Dev = "RifsxD"
const Name = "DVPL_GO CLI CONVERTER"
const Version = "4.1.0"
const Repo = "https://github.com/RifsxD/dvpl_go"
const Web = "https://rxd-mods.xyz"
const Build = "21/09/2023"
const Info = "A GUI Tool Coded In JavaScript To Convert WoTB ( Dava ) SmartDLC DVPL File Based On LZ4_HC Compression."

func main() {
	myApp := app.NewWithID("com.rxd.dvpl_go")
	myWindow := myApp.NewWindow("DVPL_GO GUI CONVERTER")

	// Load the embedded image
	iconData, _ := resources.ReadFile("resource/dvpl_go.png")
	iconResource := fyne.NewStaticResource("dvpl_go.png", iconData)
	myWindow.SetIcon(iconResource)

	config := &Config{}

	// Check if command-line arguments were provided
	if len(os.Args) > 1 {
		// Use the provided path as the initial path
		config.Path = os.Args[1]
	} else {
		// Get the current working directory and set it as the initial path
		initialPath, err := os.Getwd()
		if err != nil {
			// Handle the error, e.g., show a message to the user
			initialPath = "" // Default to an empty string if there's an error
		}
		config.Path = initialPath
	}

	compressButton := widget.NewButton("Compress", func() {
		config.Mode = "compress"
		convertFiles(myWindow, config) // Pass myWindow as a parameter
	})

	decompressButton := widget.NewButton("Decompress", func() {
		config.Mode = "decompress"
		convertFiles(myWindow, config) // Pass myWindow as a parameter
	})

	keepOriginalsCheck := widget.NewCheck("Keep Originals", func(keep bool) {
		config.KeepOriginals = keep
	})

	pathEntry := widget.NewEntry()
	pathEntry.SetText(config.Path)
	pathEntry.SetPlaceHolder("Enter directory or file path")
	pathEntry.OnChanged = func(path string) {
		config.Path = path
	}

	// Create a custom success dialog
	successDialog := dialog.NewCustom("Success", "OK", createSuccessContent(), myWindow)
	successDialog.SetDismissText("OK")

	content := container.NewVBox(
		widget.NewLabelWithStyle("DVPL_GO GUI CONVERTER • 4.1.0", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewHBox(layout.NewSpacer(), compressButton, decompressButton, layout.NewSpacer()),
		widget.NewForm(
			widget.NewFormItem("Options:", keepOriginalsCheck),
			widget.NewFormItem("Path:", pathEntry),
		),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 200))
	myWindow.ShowAndRun()
}

func showSuccessDialog(myWindow fyne.Window) {
	successDialog := dialog.NewCustom("Success", "OK", createSuccessContent(), myWindow)
	successDialog.SetDismissText("OK")
	successDialog.Show()
}

func createSuccessContent() fyne.CanvasObject {
	successLabel := widget.NewLabelWithStyle("Conversion completed successfully", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	content := container.NewVBox(
		successLabel,
	)

	return content
}

// In your convertFiles function, call showSuccessDialog when the conversion is successful.
func convertFiles(myWindow fyne.Window, config *Config) {
	err := processFiles(config.Path, config)
	if err != nil {
		dialog.NewError(err, myWindow)
	} else {
		showSuccessDialog(myWindow) // Show the custom success dialog
	}
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
