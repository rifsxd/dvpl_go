#!/bin/bash

# Function to check if the OS is Linux
is_linux() {
    if [ "$(uname -s)" == "Linux" ]; then
        return 0  # It's Linux
    else
        return 1  # It's not Linux
    fi
}

# Function to check if the OS is Windows
is_windows() {
    if [[ "$(uname -s)" == CYGWIN* || "$(uname -s)" == MINGW* ]]; then
        return 0  # It's Windows-like (Cygwin, MSYS2, Git Bash, etc.)
    else
        return 1  # It's not Windows-like
    fi
}

# Function to check if the architecture is amd64
is_amd64() {
    if [ "$(uname -m)" == "x86_64" ]; then
        return 0  # It's amd64
    else
        return 1  # It's not amd64
    fi
}

# Function to check if the architecture is 32-bit (x86)
is_x86() {
    if [ "$(uname -m)" == "i686" ]; then
        return 0  # It's x86
    else
        return 1  # It's not x86
    fi
}

# Check if the first argument is "clean" and perform cleaning if needed
if [ "$1" == "clean" ]; then

    echo ""
    echo "Cleaning build!"
    echo ""

    rm -rf "build"

    echo ""
    echo "Cleaned build!"
    echo ""
    
    mkdir "build"
fi

# Check if the OS is Linux before compiling for Linux
if is_linux && is_amd64; then
    echo ""
    echo "Compiling main package dvpl_go for linux_amd64."
    echo ""

    GOOS=linux GOARCH=amd64 go build -o "build/linux_amd64/dvpl_go"
    
    echo ""
    echo "Compiled main package dvpl_go for linux_amd64 -> build/linux_amd64/dvpl_go"
    echo ""
    
    echo ""
    echo "Compiling sub package dvpl_go_cli for linux_amd64."
    echo ""

    GOOS=linux GOARCH=amd64 go build -o "build/linux_amd64/dvpl_go_cli" "./cli"
    
    echo ""
    echo "Compiled sub package dvpl_go_cli for linux_amd64 -> build/linux_amd64/dvpl_go_cli"
    echo ""
    
    echo ""
    echo "Compiling sub package dvpl_go_gui for linux_amd64."
    echo ""

    GOOS=linux GOARCH=amd64 go build -o "build/linux_amd64/dvpl_go_gui" "./gui"
    
    echo ""
    echo "Compiled sub package dvpl_go_gui for linux_amd64 -> build/linux_amd64/dvpl_go_gui"
    echo ""
fi

# Check if the OS is Windows before compiling for Windows
if is_windows && is_amd64; then
    echo ""
    echo "Compiling main package dvpl_go for windows_amd64."
    echo ""

    GOOS=windows GOARCH=amd64 go build -o "build/windows_amd64/dvpl_go.exe"
    
    echo ""
    echo "Compiled main package dvpl_go for windows_amd64 -> build/windows_amd64/dvpl_go.exe"
    echo ""
    
    echo ""
    echo "Compiling sub package dvpl_go_cli for windows_amd64."
    echo ""

    GOOS=windows GOARCH=amd64 go build -o "build/windows_amd64/dvpl_go_cli.exe" "./cli"
    
    echo ""
    echo "Compiled sub package dvpl_go_cli for windows_amd64 -> build/windows_amd64/dvpl_go_cli.exe"
    echo ""
    
    echo ""
    echo "Compiling sub package dvpl_go_gui for windows_amd64."
    echo ""

    GOOS=windows GOARCH=amd64 go build -o "build/windows_amd64/dvpl_go_gui.exe" "./gui"
    
    echo ""
    echo "Compiled sub package dvpl_go_gui for windows_amd64 -> build/windows_amd64/dvpl_go_gui.exe"
    echo ""
fi

# Check if the OS is Linux before compiling for Linux
if is_linux && is_x86; then
    echo ""
    echo "Compiling main package dvpl_go for linux_386."
    echo ""

    GOOS=linux GOARCH=386 go build -o "build/linux_386/dvpl_go"
    
    echo ""
    echo "Compiled main package dvpl_go for linux_386 -> build/linux_386/dvpl_go"
    echo ""
    
    echo ""
    echo "Compiling sub package dvpl_go_cli for linux_386."
    echo ""

    GOOS=linux GOARCH=386 go build -o "build/linux_386/dvpl_go_cli" "./cli"
    
    echo ""
    echo "Compiled sub package dvpl_go_cli for linux_386 -> build/linux_386/dvpl_go_cli"
    echo ""
    
    echo ""
    echo "Compiling sub package dvpl_go_gui for linux_386."
    echo ""

    GOOS=linux GOARCH=386 go build -o "build/linux_386/dvpl_go_gui" "./gui"
    
    echo ""
    echo "Compiled sub package dvpl_go_gui for linux_386 -> build/linux_3864/dvpl_go_gui"
    echo ""
fi

# Check if the OS is Windows before compiling for Windows
if is_windows && is_x86; then
    echo ""
    echo "Compiling main package dvpl_go for windows_386."
    echo ""

    GOOS=windows GOARCH=386 go build -o "build/windows_386/dvpl_go.exe"
    
    echo ""
    echo "Compiled main package dvpl_go for windows_386 -> build/windows_386/dvpl_go.exe"
    echo ""
    
    echo ""
    echo "Compiling sub package dvpl_go_cli for windows_386."
    echo ""

    GOOS=windows GOARCH=386 go build -o "build/windows_386/dvpl_go_cli.exe" "./cli"
    
    echo ""
    echo "Compiled sub package dvpl_go_cli for windows_386 -> build/windows_386/dvpl_go_cli.exe"
    echo ""
    
    echo ""
    echo "Compiling sub package dvpl_go_gui for windows_386."
    echo ""

    GOOS=windows GOARCH=386 go build -o "build/windows_386/dvpl_go_gui.exe" "./gui"
    
    echo ""
    echo "Compiled sub package dvpl_go_gui for windows_386 -> build/windows_386/dvpl_go_gui.exe"
    echo ""
fi

# Pack win_64 builds
if is_windows && is_amd64; then

zip -r build/windows_amd64.zip build/windows_amd64/ 

fi

# Pack linux_64 builds
if is_linux && is_amd64; then

zip -r build/linux_amd64.zip build/linux_amd64/

fi

# Pack win_32 builds
if is_windows && is_x86; then

zip -r build/windows_386.zip build/windows_386/

fi

# Pack linux_32 builds
if is_linux && is_x86; then

zip -r build/linux_386.zip build/linux_386/

fi


# Cross compiling for windows from linux
# GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build


