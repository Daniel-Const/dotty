package main


// Represents a Dot file
type Dot struct {
    // label       string  // A label for the dot.file
    // description string  // A short description (optional)
    path        string  // Current path of the file
    deployPath  string  // Where we need to put this file (From the dotty.map)
}
