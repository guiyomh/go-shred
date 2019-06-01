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

package shred

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
)

// Shred is an object containing all choices of the user
type Shred struct {
	Iteration int64
	Zero      bool
	Remove    bool
	Jobs      int8
}

// New create a Shred structure
func New(iteration int64, zero, remove bool, jobs int8) Shred {
	if jobs == 0 {
		jobs = 3
	}
	return Shred{
		Iteration: iteration,
		Zero:      zero,
		Remove:    remove,
		Jobs:      jobs,
	}
}

// Path shreds all files in the location of path recursively.
// If zero is set to true overrides files with 0 after shredding.
// If remove is set to true files will be deleted after shredding and zero.
// When a file is shredded its content is NOT recoverable so !!USE WITH CAUTION!!
func (s Shred) Path(path string) error {
	stats, err := os.Stat(path)
	if err != nil {
		return err
	} else if stats.IsDir() {
		fmt.Printf("Listing the  directory: %s\n", path)
		err := s.dir(path)
		if err != nil {
			return err
		}
	} else {
		err := s.file(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Shred) dir(path string) error {
	jobs := make(chan string, 50)
	done := make(chan int8)
	for i := int8(1); i <= s.Jobs; i++ {
		go func(id int8, jobs <-chan string, done chan<- int8) {
			for path := range jobs {
				fmt.Printf("Job %d - running on %s ...\n", id, path)
				s.file(path)
			}
			done <- 1
		}(i, jobs, done)
	}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("Reading path : %s\n", path)
		if info.IsDir() == false && len(path) > 0 {
			jobs <- path
		}
		return nil
	})
	close(jobs)
	currentDone := int8(0)
	for d := range done {
		currentDone += d
		if currentDone == s.Jobs {
			close(done)
			break
		}
	}
	return err
}

func (s Shred) file(path string) error {
	fmt.Printf("shredding the file '%s'\n", path)
	fileinfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	size := fileinfo.Size()
	err = s.writeRandom(path, size)
	if err != nil {
		return err
	}
	err = s.writeZeros(path, size)
	if err != nil {
		return err
	}
	if s.Remove {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// writeRandom overwrites a File with random stuff.
// s.Iteration specifies how many times the file should be overwritten
func (s Shred) writeRandom(path string, size int64) error {
	var i int64
	for i = 0; i < s.Iteration; i++ {
		file, err := os.OpenFile(path, os.O_RDWR, 0)
		defer file.Close()
		if err != nil {
			return err
		}
		offset, err := file.Seek(0, 0)
		if err != nil {
			return err
		}
		buff := make([]byte, size)
		rand.Read(buff)
		file.WriteAt(buff, offset)
		file.Close()
	}
	return nil
}

// writeZeros overwrites a File with zeros if s.Zero == true
func (s Shred) writeZeros(path string, size int64) error {
	if s.Zero == false {
		return nil
	}
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	defer file.Close()
	if err != nil {
		return err
	}

	offset, err := file.Seek(0, 0)
	if err != nil {
		return err
	}
	buff := make([]byte, size)
	file.WriteAt(buff, offset)
	return nil
}
