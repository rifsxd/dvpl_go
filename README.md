# Dvpl_Go Cli & Gui Converter
- A Cli & Gui Tool Coded In Golang To Convert WoTB ( Dava ) SmartDLC DVPL File Based On LZ4_HC Compression.

 ![Demo](img/dvplgo-demo.gif)


Usage :

  - dvpl_go [-mode] [-keep-originals] [-path]

    - mode can be one of the following:

        compress: compresses files into dvpl.
        decompress: decompresses dvpl files into standard files.
		gui: opens the graphical user interface window.
        help: show this help message.

	- flags can be one of the following:

    	-keep-originals flag keeps the original files after compression/decompression.
		-path specifies the directory/files path to process. Default is the current directory.

	- usage can be one of the following examples:

		```
		$ dvpl_go -mode gui
		```
		```
		$ dvpl_go -mode help
		```
		```
		$ dvpl_go -mode decompress -path /path/to/decompress/compress
		```
		```
		$ dvpl_go -mode compress -path /path/to/decompress/compress
		```
		```
		$ dvpl_go -mode decompress -keep-originals -path /path/to/decompress/compress
		```
		```
		$ dvpl_go -mode compress -keep-originals -path /path/to/decompress/compress
		```
		```
		$ dvpl_go -mode decompress -path /path/to/decompress/compress.yaml.dvpl
		```
		```
		$ dvpl_go -mode compress -path /path/to/decompress/compress.yaml
		```
		```
		$ dvpl_go -mode decompress -keep-originals -path /path/to/decompress/compress.yaml.dvpl
		```
		```
		$ dvpl_go -mode dcompress -keep-originals -path /path/to/decompress/compress.yaml
		```


How to build :

- go 1.20+ required!

```
$ git clone https://github.com/RifsxD/dvpl-go.git
```

```
$ cd dvpl-go/src/
```

```
$ go build
```
