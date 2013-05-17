package main

import (
    "os"
    "log"
    "fmt"
    "flag"
    "time"
    "sort"
    "bufio"
    "regexp"
    "strings"
    "strconv"
)

const (
    HISTORY_SIZE = 600
    AGING_CONSTANT = 86400
    FIELD_SEP = "\x00"
    BACKUP_SUFFIX = ".bak"
)

var (
    now = time.Now().Unix()
    historySize int64
    agingConstant int64
)

type Data struct {
    path string
    hits, atime int64
}

type ByFrecency []Data

func (b ByFrecency) Len() int {
    return len(b)
}

func (b ByFrecency) Swap(i, j int) {
    b[i], b[j] = b[j], b[i]
}

func(b ByFrecency) Less(i, j int) bool {
    return Score(b[i].hits, now - b[i].atime) < Score(b[j].hits, now - b[j].atime)
}

func Score(hits int64, age int64) (float64) {
    return float64(hits) * float64(agingConstant) / float64(agingConstant + age)
}

func check(e error) {
    if (e != nil) {
        panic(e)
    }
}

func ReadData(r *bufio.Reader) (Data, error) {
    line, err := r.ReadString('\n')
    if (err != nil) {
        return Data{}, err
    }
    line = line[:len(line) - 1]
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
    hs := os.Getenv("Z_HISTORY_SIZE")
    if historySize, _ = strconv.ParseInt(hs, 10, 64); historySize < 1 {
        historySize = HISTORY_SIZE
    }
    results := make([]Data, 0, historySize)
    flag.Int64Var(&agingConstant, "g", AGING_CONSTANT, "Set the value of the aging constant")
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
        pattern, err := regexp.Compile("(?i)" + sPattern)
        check(err)
        for d, err := ReadData(bf); err == nil; d, err = ReadData(bf) {
            if pattern.MatchString(d.path) {
                results = append(results, d)
            }
        }
        sort.Sort(sort.Reverse(ByFrecency(results)))
        for _, d := range results {
            fmt.Printf("%v\n", d.path)
        }
    } else {
        var index int64 = -1
        var cur int64
        for d, err := ReadData(bf); err == nil; d, err = ReadData(bf) {
            if (index < 0 && pathFlag == d.path) {
                index = cur
            }
            results = append(results, d)
            cur++
        }
        if len(*addFlag) != 0 {
            if (index < 0) {
                results = append(results, Data{pathFlag, 1, now})
            } else {
                results[index] = Data{pathFlag, results[index].hits + 1, now}
            }
            sort.Sort(sort.Reverse(ByFrecency(results)))
        } else if len(*deleteFlag) != 0 {
            if (index < 0) {
                log.Printf("Item is missing: '%v'.", pathFlag)
            } else {
                results = append(results[:index], results[index + 1:]...)
            }
        }
        dataFileBackup := dataFile + BACKUP_SUFFIX
        fobj, err := os.Create(dataFileBackup)
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
        err = os.Rename(dataFileBackup, dataFile)
        check(err)
    }
}
