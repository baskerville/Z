package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHistorySize   = 600
	fieldSep             = "\x00"
)

var (
	now           = time.Now().Unix()
	historySize   int64
)

type Data struct {
	path        string
	hits, atime int64
}

type Datae []Data

func (d Datae) Len() int {
	return len(d)
}

func (d Datae) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d Datae) Less(i, j int) bool {
	return Score(d[i].hits, now-d[i].atime) < Score(d[j].hits, now-d[j].atime)
}

func Score(hits int64, age int64) float64 {
	return float64(hits) / (0.25 + float64(age) * 3e-16)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadData(r *bufio.Reader) (Data, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return Data{}, err
	}
	line = line[:len(line)-1]
	tok := strings.Split(line, fieldSep)
	atime, err := strconv.ParseInt(tok[0], 10, 64)
	check(err)
	hits, err := strconv.ParseInt(tok[1], 10, 64)
	check(err)
	path := tok[2]
	return Data{path, hits, atime}, nil
}

func main() {
	var dataFile string
	var pathFlag string
	if dataFile = os.Getenv("Z_DATA_FILE"); len(dataFile) == 0 {
		dataFile = os.Getenv("HOME") + string(os.PathSeparator) + ".z"
	}
	if historySize, _ = strconv.ParseInt(os.Getenv("Z_HISTORY_SIZE"), 10, 64); historySize < 1 {
		historySize = defaultHistorySize
	}
	results := make(Datae, 0, historySize)
	addFlag := flag.String("a", "", "Add the given item to the data file")
	deleteFlag := flag.String("d", "", "Delete the given item from the data file")
	inputFlag := flag.String("i", dataFile, "Use the given file as data file")
	flag.Parse()
	dataFile = *inputFlag
	if len(*addFlag) != 0 {
		pathFlag = *addFlag
	} else if len(*deleteFlag) != 0 {
		pathFlag = *deleteFlag
	}
	fobj, err := os.Open(dataFile)
	check(err)
	defer fobj.Close()
	var bf = bufio.NewReader(fobj)
	if len(pathFlag) == 0 {
		var sPattern string
		if sPattern = strings.Join(flag.Args(), ".*"); len(sPattern) == 0 {
			sPattern = ".*"
		}
		reFlags := "(?i)"
		if sPattern != strings.ToLower(sPattern) {
			reFlags = ""
		}
		pattern, err := regexp.Compile(reFlags + sPattern)
		check(err)
		for d, err := ReadData(bf); err == nil; d, err = ReadData(bf) {
			if pattern.MatchString(d.path) {
				results = append(results, d)
			}
		}
		sort.Sort(sort.Reverse(Datae(results)))
		for _, d := range results {
			fmt.Printf("%v\n", d.path)
		}
	} else {
		var index int64 = -1
		var cur int64
		for d, err := ReadData(bf); err == nil; d, err = ReadData(bf) {
			if index < 0 && pathFlag == d.path {
				index = cur
			}
			results = append(results, d)
			cur++
		}
		if len(*addFlag) != 0 {
			if index < 0 {
				results = append(results, Data{pathFlag, 1, now})
			} else {
				results[index] = Data{pathFlag, results[index].hits + 1, now}
			}
			if int64(len(results)) > historySize {
				sort.Sort(Datae(results))
				results = results[1:]
			}
		} else if len(*deleteFlag) != 0 {
			if index < 0 {
				log.Printf("Item is missing: '%v'.", pathFlag)
			} else {
				results = append(results[:index], results[index+1:]...)
			}
		}
		fobj, err := ioutil.TempFile(filepath.Dir(dataFile), filepath.Base(dataFile))
		check(err)
		defer fobj.Close()
		bf := bufio.NewWriter(fobj)
		for _, d := range results {
			_, err := bf.WriteString(fmt.Sprintf("%v%v%v%v%v\n", d.atime, fieldSep, d.hits, fieldSep, d.path))
			check(err)
		}
		err = bf.Flush()
		check(err)
		err = os.Rename(fobj.Name(), dataFile)
		check(err)
	}
}
