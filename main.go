// Copyright (C) 2019 Guillaume Camus
//
// This file is part of go-shred.
//
// go-shred is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-shred is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-shred.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	flags "github.com/jessevdk/go-flags"
	"githuh.com/guiyomh/go-shred/pkg/shred"
)

const (
	Name    = "go-shred"
	Author  = "Guillaume CAMUS"
	Version = "0.1.0"
)

var opts struct {
	Iteration   int64 `short:"n" long:"iteration"  description:"overwrite N times instead of the default" default:"3"`
	Remove      bool  `          long:"remove"     description:"deallocate and remove file after overwriting"`
	Zero        bool  `short:"z" long:"zero"       description:"add a final overwrite with zeros to hide shredding"`
	Jobs        int8  `short:"j" long:"jobs"       description:"number of job to apply the shredding" default:"16"`
	Verbose     bool  `short:"v" long:"verbose"    description:"verbose mode"`
	ShowVersion bool  `short:"V" long:"version"    description:"show version and exit"`
	ShowHelp    bool  `short:"h" long:"help"       description:"show this help message"`
	Args        struct {
		Path []string `description:"List of directory/files to shred"`
	} `positional-args:"yes" required:"yes"`
}

func initArgParser() []string {
	var err error
	argparser := flags.NewParser(&opts, flags.PassDoubleDash)
	args, err := argparser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if opts.ShowVersion {
		fmt.Printf("%s version %s\n", Name, Version)
		fmt.Printf("Copyright (C) 2019 %s\n", Author)
		os.Exit(0)
	}
	if opts.ShowHelp {
		argparser.WriteHelp(os.Stdout)
		os.Exit(1)
	}
	if len(opts.Args.Path) == 0 {
		fmt.Printf("Error: You must specifie a Path file or directory\n\n")
		argparser.WriteHelp(os.Stdout)
		os.Exit(1)
	}
	return args
}

func main() {
	start := time.Now()
	initArgParser()
	shred := shred.New(opts.Iteration, opts.Zero, opts.Remove, opts.Jobs)
	var wg sync.WaitGroup
	for _, path := range opts.Args.Path {
		fmt.Printf("Start the shredding of %s\n", path)
		wg.Add(1)
		go shredFn(shred, path, &wg)
	}
	wg.Wait()
	elapsed := time.Now().Sub(start)
	fmt.Printf("Shredding finish in %s\n", elapsed)
}

func shredFn(s shred.Shred, path string, wg *sync.WaitGroup) {
	defer wg.Done()
	s.Path(path)
}
