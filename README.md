# dvpl-go cli converter
 A CLI Tool Coded In JavaScript To Convert WoTB ( Dava ) SmartDLC DVPL File Based On LZ4 Comrpession.

 ![Demo](img/dvplgo-demo.gif)

```

Usage :

  dvpl [-mode] [-keep-originals] [-path]

    • mode can be one of the following:

        compress: compresses files into dvpl.
        decompress: decompresses dvpl files into standard files.
        help: show this help message.

	• flags can be one of the following:

    	-keep-originals flag keeps the original files after compression/decompression.
    	-path specifies the directory path to process. Default is the current directory.

	• usage can be one of the following examples:

$ dvplgo -mode decompress -path /path/to/decompress/compress
		
$ dvplgo -mode compress -path /path/to/decompress/compress
		
$ dvplgo -mode decompress -keep-originals -path /path/to/decompress/compress
		
$ dvplgo -mode compress -keep-originals -path /path/to/decompress/compress
		
$ dvplgo -mode decompress -path /path/to/decompress/compress.yaml.dvpl
		
$ dvplgo -mode compress -path /path/to/decompress/compress.yaml
		
$ dvplgo -mode decompress -keep-originals -path /path/to/decompress/compress.yaml.dvpl
		
$ dvplgo -mode dcompress -keep-originals -path /path/to/decompress/compress.yaml


```

Building :

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