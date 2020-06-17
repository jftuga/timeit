# timeit
A cross-platform CLI tool used to time the duration of the given command 

## Usage

```
Usage: timeit [cmd] [args...]
You may need to surround args within double-quotes

Examples:
    timeit wget https://example.com/file.tar.gz
    timeit gzip -t "file with spaces.gz"

For built-in Windows 'cmd' commands:
    timeit cmd /c "dir c:\ /s/b > list.txt"
    timeit cmd /c dir /s "c:\Program Files"

```

## Examples
```
c:\>timeit.exe sleep 2
2.04042s
```

```
pi@pi:~ $ timeit md5sum linux.iso
4.1329682s
```
