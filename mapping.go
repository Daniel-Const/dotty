package main

/*
 * Functions for parsing dotty.map files
 */

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func ReadDotMap(profileDirPath string)  []*Dot {
    var dots []*Dot

    // Read .map file and create a slice of dots
    file, err := os.Open(profileDirPath + "/dotty.map")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.Split(line, ":")
        if len(parts) < 2 {
            log.Fatal("Err: Map file incorrectly formatted")
        }
        fileName := strings.Trim(parts[0], " ")
        toPath   := strings.Trim(parts[1], " ")

        dot := Dot{profileDirPath+"/"+fileName, toPath}
        dots = append(dots, &dot)
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    return dots
}
