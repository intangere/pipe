# pipe
Generic pipe "operator" for Go with fine-grained control.

# What does this do?    
Effectively the same functional composition the `|>` operator in Elixir allows for plus a little more. You can pipe any amount of outputs of one function as inputs to another (assuming the function expects that many inputs), automatic error handling, partial execution, examine the current state of the pipeline, execute further functions ontop of an existing pipeline, and deferred execution. Additionally, the result of the execution of a pipeline is unwrapped into a concrete type thanks to generics.

# Install: 
```` 
go get github.com/intangere/pipe 
````    

Usage (for more see *examples/example.go*): 
````go

    // this is an awful example, but it demonstrates a few concepts

    // Instead of doing this
    example := "Hello World!"
    example = strings.ToLower(example)
    splits := strings.Split(example, " ")
    new := "Bye " + splits[1]
    log.Println("Result:", new)

    // With pipe you can do
    res, err := Pipe[string]("Hello World!").
        Flow(strings.ToLower).
        Flow(strings.Split, " ").
        Flow(func (parts []string) string {
                return "Bye " + parts[1]
        }).
        Unwrap()

    log.Println("Result:", res, "Error:", err)

    // Output: Byte world!

````
Credits:   
https://github.com/aslrousta/pipe - reflect code copied nearly 1-1 from here
