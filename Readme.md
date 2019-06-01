<!--
 Copyright (C) 2019 Guillaume Camus

 This file is part of go-shred.

 go-shred is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 go-shred is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with go-shred.  If not, see <http://www.gnu.org/licenses/>.
-->

[![GoDoc](https://godoc.org/github.com/guiyomh/go-shred?status.svg)](https://godoc.org/github.com/guiyomh/go-shred)
[![License](https://img.shields.io/github/license/guiyomh/go-shred.svg)](https://www.gnu.org/licenses/gpl-3.0.en.html)
# go-shred
This is a golang command line and a libary to mimic the functionallity of the linux ```shred``` command

## Command line usage
```bash
./go-shred -h
  go-shred Path...

Application Options:
  -n, --iteration= overwrite N times instead of the default (default: 3)
      --remove     deallocate and remove file after overwriting
  -z, --zero       add a final overwrite with zeros to hide shredding
  -j, --jobs=      number of job to apply the shredding (default: 16)
  -v, --verbose    verbose mode
  -V, --version    show version and exit
  -h, --help       show this help message

Arguments:
  Path:            List of directory/files to shred
```

## library Usage
```golang
package main
import (
    "githuh.com/guiyomh/go-shred/pkg/shred"
)
type conf struct {
    Iteration   int64
    Remove      bool
    Zero        bool
    Jobs        int8
}

func main(){
    ops := conf{
        Iteration: 3,
        Remove: true,
        Zero: false,
        Jobs: 16,
    }
    shred := shred.New(opts.Iteration, opts.Zero, opts.Remove, opts.Jobs)
    shred.Path("filename")
}
```