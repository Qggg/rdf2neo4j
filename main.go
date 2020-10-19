package main

import (
    "os"
    "flag"
    "io"
    "bufio"
    "fmt"
    "time"
    "strconv"
    "encoding/csv"
    "strings"
)

var rdfPath = flag.String("path", "", "Specify rdf data path")

const SEED = 0xc70f6907

const CAP = 100000

func main() {
    flag.Parse()
  
    if err := Read(*rdfPath); err != nil {
        panic(err)
    }
}

func Read(rdfPath string) error {
    srcFile, err := os.Open(rdfPath)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    var vFilePath = "./vertex.csv"
    vFile, err := os.Create(vFilePath)
    if err != nil {
        return err
    }
    defer vFile.Close()

    var eFilePath = "./edge.csv"
    eFile, err := os.Create(eFilePath)
    if err != nil {
        return err
    }
    defer eFile.Close()

    reader := csv.NewReader(bufio.NewReader(srcFile))
    vWriter := csv.NewWriter(bufio.NewWriter(vFile))
    eWriter := csv.NewWriter(bufio.NewWriter(eFile))

    vertexHeadline := []string{":ID", "name", ":LABEL"}
    edgeHeadline := []string{":START_ID", "name", ":END_ID", ":TYPE"}
    vWriter.Write(vertexHeadline)
    eWriter.Write(edgeHeadline)
    defer vWriter.Flush()
    defer eWriter.Flush()

    now := time.Now()

    defer func() {
        fmt.Printf("Finish convert rdf data, consume time: %.2fs\n", time.Since(now).Seconds())
    }()

    lineNum, numErrorLines := 0, 0

    lVRecord := make([]string, 3)
    rVRecord := make([]string, 3)
    eRecord := make([]string, 4)
    exists := make(map[int64]bool)   

    // re := regexp.MustCompile(`\r?\n`)
    for {
        if lineNum % 100000 == 0 {
            fmt.Printf("hava read lines: %d\n", lineNum)
        }

        line, err := reader.Read()
        if err == io.EOF {
            fmt.Printf("totalLines: %d, errorLines: %d\n", lineNum, numErrorLines)
            break
        }

        lineNum++

        if err != nil {
            numErrorLines++
            continue
        }

        if len(line) != 3 {
            numErrorLines++
            continue
        }

        if lineNum % 100000 == 0 {
            fmt.Printf("seg1=%s,seg2=%s,seg3=%s", line[0], line[1], line[2])
        }
        
        if len(exists) >= CAP {
            exists = make(map[int64]bool)
        }

        hashVal := MurmurHash64A([]byte(line[0]), SEED)
        vid := int64(hashVal)
        lVRecord[0] = strconv.FormatInt(vid, 10)
        // line[0] = strings.Replace(line[2], " ", "", -1)
        line[0] = strings.Replace(line[2], "\t", "", -1)
        line[0] = strings.Replace(line[2], "\n", "", -1)
        line[0] = strings.Replace(line[2], "\r", "", -1)
        lVRecord[1] = line[0]
        lVRecord[2] = "ENTITY"
        if _, ok := exists[vid]; !ok {
            exists[vid] = true
            vWriter.Write(lVRecord)
        }

        hashVal = MurmurHash64A([]byte(line[2]), SEED)
        vid = int64(hashVal)
        rVRecord[0] = strconv.FormatInt(vid, 10)
        // line[2] = strings.Replace(line[2], " ", "", -1)
        line[2] = strings.Replace(line[2], "\t", "", -1)
        line[2] = strings.Replace(line[2], "\n", "", -1)
        line[2] = strings.Replace(line[2], "\r", "", -1)
        rVRecord[1] = line[2]
        rVRecord[2] = "ENTITY"
        if _, ok := exists[vid]; !ok {
            exists[vid] = true
            vWriter.Write(rVRecord)
        }
    
        eRecord[0] = lVRecord[0]
        eRecord[2] = rVRecord[0]
        eRecord[1] = line[1]
        eRecord[3] = line[1]
        eWriter.Write(eRecord)
    }

    return nil
}
