# pipe
Generic pipe "operator" for Go with fine-grained control.

# What does this do?    
Effectively the same functional composition the `|>` operator in Elixir allows for plus a little more. You can pipe any amount of outputs of one function as inputs to another (assuming the function expects that many inputs), automatic error handling, partial execution, examine the current state of the pipeline, execute further functions ontop of an existing pipeline, and deferred execution. Additionally, the result of the execution of a pipeline is unwrapped into a concrete type thanks to generics.

# Install: 
```` 
go get github.com/intangere/pipe 
````    

Usage (with execution using reflection)  (for more see *examples/example.go*): 
````go

    import . "github.com/intangere/pipe"
    // this is an awful example, but it demonstrates a few concepts

    // Instead of doing this
    example := "Hello World!"
    example = strings.ToLower(example)
    splits := strings.Split(example, " ")
    new := "Bye " + splits[1]
    log.Println("Result:", new)
    // Output: Bye world!

    // With pipe you can do
    res, err := Pipe[string]("Hello World!").
        Flow(strings.ToLower).
        Flow(strings.Split, " ").
        Flow(func (parts []string) string {
                return "Bye " + parts[1]
        }).
        Unwrap()

    log.Println("Result:", res, "Error:", err)
    // Output: Bye world!

    // Note: The generic type parameter for Pipe[T] is the type of the result from your pipeline

    // or with deferred execution
    res, err = Pipe[string]("Hello World!").
        Next(strings.ToLower).
        Next(strings.Split, " ").
        Next(func (parts []string) string {
                return "Bye " + parts[1]
        }).
        Do().
        Unwrap()

    log.Println("Result:", res, "Error:", err)
    // Output: Bye world!

    //reusable pipe
    p = Pipe[string]().
        Next(func (s string) string {
           return "With argument: " + s
        }).
        Next(log.Println)

    for i := 0; i < 3; i++ {
        // no result
        _, err = p.Result(fmt.Sprintf("%d", i))
    }

    // will print:
    // With argument: 0
    // With argument: 1
    // With argument: 2

````

Usage (without reflection, constrained to 1 input/1 output/1 type, no magic)  (for more see *examples/simple.go*): 
````go

    import . "github.com/intangere/pipe/simple"

    func main() {
        res, err := Pipe[string]("Hello World!\n  ").
                Flow(strings.ToLower).
                Flow(strings.Title).
                Flow(strings.TrimSpace).
                Result()

        log.Println("Res: ", res, "Error:", err)
    }
```
Credits:   
https://github.com/aslrousta/pipe - reflect code copied nearly 1-1 from here
