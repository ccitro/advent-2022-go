package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type file struct {
	name string
	size int
}

type directory struct {
	name    string
	parent  *directory
	files   []*file
	subdirs []*directory
}

func getDirSize(d *directory) int {
	size := 0
	for _, f := range d.files {
		size += f.size
	}
	for _, sd := range d.subdirs {
		size += getDirSize(sd)
	}
	return size
}

func printDirectory(d *directory, indent int) {
	print(strings.Repeat(" ", indent))
	println("-", d.name, "(dir)")
	indent += 2
	for _, f := range d.files {
		print(strings.Repeat(" ", indent))
		println("-", f.name, "(file, size="+fmt.Sprintf("%d", f.size)+")")
	}
	for _, sd := range d.subdirs {
		printDirectory(sd, indent)
	}
}

func (d *directory) print() {
	printDirectory(d, 0)
}

func parseFilesystem(f *os.File) directory {
	rootdir := directory{name: "/"}
	cwd := &rootdir

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "$") {
			command := line[2:]
			if command == "cd /" {
				cwd = &rootdir
			} else if command == "cd .." {
				cwd = cwd.parent
			} else if strings.HasPrefix(command, "cd ") {
				new_dir_name := command[3:]
				for _, sd := range cwd.subdirs {
					if sd.name == new_dir_name {
						cwd = sd
						break
					}
				}
			}
		} else if strings.HasPrefix(line, "dir") {
			dir_name := line[4:]
			new_dir := directory{name: dir_name, parent: cwd}
			cwd.subdirs = append(cwd.subdirs, &new_dir)
		} else {
			split := strings.Split(line, " ")
			size, _ := strconv.Atoi(split[0])
			name := split[1]
			new_file := file{name: name, size: size}
			cwd.files = append(cwd.files, &new_file)
		}
	}

	return rootdir
}

func findDirsUnderSize(d *directory, max_size int) []*directory {
	dirs := []*directory{}
	s := getDirSize(d)
	if s <= max_size {
		dirs = append(dirs, d)
	}

	for _, sd := range d.subdirs {
		dirs = append(dirs, findDirsUnderSize(sd, max_size)...)
	}

	return dirs
}

func part1(file *os.File) {
	dir := parseFilesystem(file)

	dirs := findDirsUnderSize(&dir, 100000)
	accum := 0
	for _, d := range dirs {
		accum += getDirSize(d)
	}
	println(accum)
}

func part2(file *os.File) {
	dir := parseFilesystem(file)

	fs_size := 70000000
	space_req := 30000000
	space_used := getDirSize(&dir)
	space_avail := fs_size - space_used
	delete_size_req := space_req - space_avail

	delete_size := space_used

	dirs_to_search := []*directory{&dir}
	for len(dirs_to_search) > 0 {
		d := dirs_to_search[0]
		dirs_to_search = dirs_to_search[1:]

		s := getDirSize(d)
		if s < delete_size_req {
			continue
		}
		if s < delete_size {
			delete_size = s
		}

		dirs_to_search = append(dirs_to_search, d.subdirs...)
	}

	println(delete_size)
}

func main() {
	filename := "input.txt"
	method := part1
	for _, v := range os.Args {
		if v == "part2" || v == "2" {
			method = part2
		}
		if strings.HasSuffix(v, ".txt") {
			filename = v
		}
	}

	file, _ := os.Open(filename)
	defer file.Close()
	method(file)
}
