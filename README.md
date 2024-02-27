# pipe
Generic pipe "operator" for Go with fine-grained control.

# What does this do?    
Effectively the same functional composition the `|>` operator in Elixir allows for plus a little more. You can pipe any amount of outputs of one function as inputs to another (assuming the function expects that many inputs), automatic error handling, partial execution, examine the current state of the pipeline, and, execute further functions ontop of an existing pipeline, and deferred execution. Additionally, the result of the execution of a pipeline is unwrapped into a concrete type thanks to generics.

# Install: 
```` 
go get github.com/intangere/pipe 
````    
Usage (for more see *examples/example.go*): 
````go
    // deferred execution
    res, err = Pipe[string]("Hello").
	Next(strings.ToLower).
	Next(func (s string) string {
		return s + " world!"
	}).
	Next(strings.Title).
    Do().
	Unwrap()

    log.Println("Result:", res, "Error:", err)
    // Output: Hello world!

    // execute as the pipe as it is built
    res, err = Pipe[string]("Hello").
        Flow(strings.ToLower).
        Flow(func (s string) string {
                return s + " world!"
        }).
        Flow(strings.Title).
        Unwrap()

    log.Println("Result:", res, "Error:", err)
    // Output: Hello world!

````
Credits:   
https://github.com/aslrousta/pipe - reflect code copied nearly 1-1 from here
