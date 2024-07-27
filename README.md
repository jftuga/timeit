# timeit
A cross-platform CLI tool used to time the duration of the given command

## Installation

* [Binaries](https://github.com/jftuga/timeit/releases) are available for Windows, MacOS and Linux
* Homebrew installation for MacOS and Linux:
* * `brew tap jftuga/homebrew-tap; brew update; brew install jftuga/tap/timeit`

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

Run 'timeit _start' to create (or overwrite) this file containing the current time: .timeit.start.tmp
Run 'timeit _end' to read (and then delete) that file.  The elapsed time will then be displayed.
This can be useful for timing multiple, long-running commands.

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

```
c:\>timeit _start
c:\>long_running_process 1
c:\>long_running_process 2
c:\>long_running_process 3
c:\timeit _end
6h0m57.2788758s
```
