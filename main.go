package main

import (
    "fmt"
    "github.com/jasonlvhit/gocron"
    "io/ioutil"
    "log"
    "path"
    "time"
)

func dirRotate(dir string, months int) {
    now := time.Now()
    fmt.Println("now:", now)
    then := now.AddDate(0, -months, 0)
    fmt.Println("then:", then)
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
                // TODO remove
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
                // TODO remove
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
                // TODO remove
            }
        }
    }
}

func main() {
    gocron.Every(10).Seconds().Do(dirRotate, "/tmp/log_bak", 1)
    <-gocron.Start()
}
