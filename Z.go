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
    HISTORY_SIZE     = 600
    AGING_CONSTANT   = 86400
    FIELD_SEP        = "\x00"
    SORT_BY_FRECENCY = "frecency"
    SORT_BY_HITS     = "hits"
    SORT_BY_ATIME    = "atime"
)

var (
    now           = time.Now().Unix()
    historySize   int64
    agingConstant int64
)

type Data struct {
    path        string
    hits, atime int64
}

type Datae []Data

type ByFrecency struct {
    Datae
}

type ByAtime struct {
    Datae
}

type ByHits struct {
    Datae
}

func (d Datae) Len() int {
    return len(d)
}

func (d Datae) Swap(i, j int) {
    d[i], d[j] = d[j], d[i]
}

func (b ByFrecency) Less(i, j int) bool {
    return Score(b.Datae[i].hits, now-b.Datae[i].atime) < Score(b.Datae[j].hits, now-b.Datae[j].atime)
}

func (b ByHits) Less(i, j int) bool {
    return b.Datae[i].hits < b.Datae[j].hits
}

func (b ByAtime) Less(i, j int) bool {
    return b.Datae[i].atime < b.Datae[j].atime
}

func Score(hits int64, age int64) float64 {
    return float64(hits) * float64(agingConstant) / float64(agingConstant+age)
}

func (d Datae) Sort(method string) {
    if method == SORT_BY_FRECENCY {
        sort.Sort(sort.Reverse(ByFrecency{d}))
    } else if method == SORT_BY_HITS {
        sort.Sort(sort.Reverse(ByHits{d}))
    } else if method == SORT_BY_ATIME {
        sort.Sort(sort.Reverse(ByAtime{d}))
    } else {
        log.Fatalf("Unknown sort method: '%v'.", method)
    }
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
    tok := strings.Split(line, FIELD_SEP)
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
        historySize = HISTORY_SIZE
    }
    if agingConstant, _ = strconv.ParseInt(os.Getenv("Z_AGING_CONSTANT"), 10, 64); agingConstant < 1 {
        agingConstant = AGING_CONSTANT
    }
    results := make(Datae, 0, historySize)
    addFlag := flag.String("a", "", "Add the given item to the data file")
    deleteFlag := flag.String("d", "", "Delete the given item from the data file")
    sortFlag := flag.String("s", SORT_BY_FRECENCY, "Use the given sort method")
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
        pattern, err := regexp.Compile("(?i)" + sPattern)
        check(err)
        for d, err := ReadData(bf); err == nil; d, err = ReadData(bf) {
            if pattern.MatchString(d.path) {
                results = append(results, d)
            }
        }
        results.Sort(*sortFlag)
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
            results.Sort(*sortFlag)
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
        cur = 0
        for _, d := range results {
            if cur >= historySize {
                break
            }
            _, err := bf.WriteString(fmt.Sprintf("%v%v%v%v%v\n", d.atime, FIELD_SEP, d.hits, FIELD_SEP, d.path))
            check(err)
            cur++
        }
        err = bf.Flush()
        check(err)
        err = os.Rename(fobj.Name(), dataFile)
        check(err)
    }
}
