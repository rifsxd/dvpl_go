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

# Check if the first argument is "clean" and perform cleaning if needed
if [ "$1" == "clean" ]; then
    echo "Cleaning build!"
    rm -rf "build"
    echo "Cleaned build!"
    mkdir "build"
fi

# Check if the OS is Linux before compiling for Linux
if is_linux; then
    echo ""
    echo "Compiling main package dvpl_go for linux_amd64."
    echo ""

    go build -o "build/linux_amd64/dvpl_go"
    
    echo ""
    echo "Compiled main package dvpl_go for linux_amd64 -> build/linux_amd64/dvpl_go"
    echo ""
    
    echo ""
    echo "Compiling sub package dvpl_go_cli for linux_amd64."
    echo ""

    go build -o "build/linux_amd64/dvpl_go_cli" "./cli"
    
    echo ""
    echo "Compiled sub package dvpl_go_cli for linux_amd64 -> build/linux_amd64/dvpl_go_cli"
    echo ""
    
    echo ""
    echo "Compiling sub package dvpl_go_gui for linux_amd64."
    echo ""

    go build -o "build/linux_amd64/dvpl_go_gui" "./gui"
    
    echo ""
    echo "Compiled sub package dvpl_go_gui for linux_amd64 -> build/linux_amd64/dvpl_go_gui"
    echo ""
fi

# Check if the OS is Windows before compiling for Windows
if is_windows; then
    echo ""
    echo "Compiling main package dvpl_go for windows_amd64."
    echo ""

    go build -o "build/windows_amd64/dvpl_go.exe"
    
    echo ""
    echo "Compiled main package dvpl_go for windows_amd64 -> build/windows_amd64/dvpl_go.exe"
    echo ""
    
    echo ""
    echo "Compiling sub package dvpl_go_cli for windows_amd64."
    echo ""

    go build -o "build/windows_amd64/dvpl_go_cli.exe" "./cli"
    
    echo ""
    echo "Compiled sub package dvpl_go_cli for windows_amd64 -> build/windows_amd64/dvpl_go_cli.exe"
    echo ""
    
    echo ""
    echo "Compiling sub package dvpl_go_gui for windows_amd64."
    echo ""

    go build -o "build/windows_amd64/dvpl_go_gui.exe" "./gui"
    
    echo ""
    echo "Compiled sub package dvpl_go_gui for windows_amd64 -> build/windows_amd64/dvpl_go_gui.exe"
    echo ""
fi
