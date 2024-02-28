package main

import . "github.com/intangere/pipe"
import "strings"
import "time"
import "log"

func main() {
    t := time.Now()

    for i := 0; i < 100000; i++ {
            // instead of doing this
            example := "Hello World!"
            example = strings.ToLower(example)
            splits := strings.Split(example, " ")
            new := "Bye " + splits[1]
                new = new + " "
            //log.Println("Result:", new)
    }

    total := time.Now().Sub(t)

    concat := func (parts []string) string { 
                return "Bye " + parts[1]
    }

        t = time.Now()
            for i := 0; i < 100000; i++ {
            // with pipe you can do
            res, err := Pipe[string]("Hello World!").
                Flow(strings.ToLower).
                Flow(strings.Split, " ").
                Flow(concat).
                Unwrap()

            _ = res
                _ =  err
            //log.Println("Result:", res, "Error:", err)
        }

    total1 := time.Now().Sub(t)

    t = time.Now()

    // reusing a pipe
        p := Pipe[string]().
                Next(strings.ToLower).
                Next(strings.Split, " ").
                Next(concat)

            for i := 0; i < 100000; i++ {
            // with pipe you can do
            res, err := p.Result("Hello world!")

            _ = res
                _ =  err
            //log.Println("Result:", res, "Error:", err)
        }

    total2 := time.Now().Sub(t)

        t = time.Now()
            for i := 0; i < 100000; i++ {
            // with pipe you can do
            res, err := Pipe[string]("Hello World!").
                Next(strings.ToLower).
                Next(strings.Split, " ").
                Next(concat).
                Result()

            _ = res
                _ =  err
            //log.Println("Result:", res, "Error:", err)
        }

    total3 := time.Now().Sub(t)

    log.Println(total, total1, total2, total3, total/100000, total1/100000, total2/100000, total3/100000)
}
