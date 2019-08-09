package main

import (
    "fmt"
    "github.com/jasonlvhit/gocron"
    "gopkg.in/alecthomas/kingpin.v2"
    "io/ioutil"
    "log"
    "os"
    "path"
    "time"
)

var (
    rotateRoot = kingpin.Flag("rotate-root", "切割根目录").Envar("ROTATE_ROOT").Default("/opt/log_bak").String()
    keepMonth  = kingpin.Flag("keep-month", "保留最近几个月").Envar("KEEP_MONTH").Default("6").Int()
)

func dirRotate(dir string, months int) {
    now := time.Now()
    fmt.Println("now:", now)
    then := now.AddDate(0, -months, 0)
    fmt.Println("then:", then)
    log.Printf("begin to monitoring dir %s", dir)
    removeYearDir(dir, then)
}

// 传入目录， 目录结构类型 1:Year 2:Month 3:Day, <= dateNum
func removeYearDir(dir string, then time.Time) {
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
    }
    for _, yearFile := range files {
        if yearFile.IsDir() {
            yearTime, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("%s-01-01", yearFile.Name()), time.Local)
            fmt.Printf("year time is %v", yearTime)
            if err != nil {
                log.Printf("Dir %s is not a year dir", yearFile.Name())
                continue
            }
            // 年份<
            if yearTime.Year() < then.Year() {
                log.Printf("remove year dir %s", yearFile.Name())
                if err := os.RemoveAll(path.Join(dir, yearFile.Name())); err != nil {
                    log.Printf("remove year dir %s error", yearFile.Name())
                }
            } else if yearTime.Year() == then.Year() {
                // 对比月份
                go removeMonthDir(path.Join(dir, yearFile.Name()), then)
            } else {
                continue
            }
        }
    }
}

func removeMonthDir(dir string, then time.Time) {
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
    }
    for _, monthFile := range files {
        if monthFile.IsDir() {
            monthTime, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("2016-%s-01", monthFile.Name()), time.Local)
            fmt.Printf("month time is %v", monthTime)
            if err != nil {
                log.Printf("Dir %s is not a month dir", monthFile.Name())
                continue
            }
            // 月份<
            if monthTime.Month() < then.Month() {
                log.Printf("remove month dir %s", monthFile.Name())
                if err := os.RemoveAll(path.Join(dir, monthFile.Name())); err != nil {
                    log.Printf("remove month dir %s error", monthFile.Name())
                }
            } else if monthTime.Month() == then.Month() {
                // 对比月份
                go removeDayDir(path.Join(dir, monthFile.Name()), then)
            } else {
                continue
            }
        }
    }
}

func removeDayDir(dir string, then time.Time) {
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
    }
    for _, dayFile := range files {
        if dayFile.IsDir() {
            dayTime, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("2016-08-%s", dayFile.Name()), time.Local)
            fmt.Printf("day time is %v", dayTime)
            if err != nil {
                log.Printf("Dir %s is not a day dir", dayFile.Name())
                continue
            }
            // 日期<
            if dayTime.Day() < then.Day() {
                log.Printf("remove day dir %s", dayFile.Name())
                if err := os.RemoveAll(path.Join(dir, dayFile.Name())); err != nil {
                    log.Printf("remove day dir %s error", dayFile.Name())
                }
            }
        }
    }
}

func main() {
    kingpin.Parse()
    dirRotate(*rotateRoot, *keepMonth)
    gocron.Every(1).Day().Do(dirRotate, *rotateRoot, *keepMonth)
    <-gocron.Start()
}
