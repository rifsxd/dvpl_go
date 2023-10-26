package cli_gui

import (
	"embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

//go:embed resource/dvpl_go.png
var resources embed.FS

func Gui() {
	myApp := app.NewWithID("com.rxd.dvpl_go")
	myWindow := myApp.NewWindow("DVPL_GO GUI CONVERTER")

	// Load the embedded image
	iconData, _ := resources.ReadFile("resource/dvpl_go.png")
	iconResource := fyne.NewStaticResource("dvpl_go.png", iconData)
	myWindow.SetIcon(iconResource)

	config := &Config{}

	/* Check if command-line arguments were provided
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
	} */

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
		widget.NewLabelWithStyle("DVPL_GO GUI CONVERTER • 4.2.0", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
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
