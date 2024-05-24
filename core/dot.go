package core 


// Represents a Dot file
type Dot struct {
    // label       string  // A label for the dot.file
    // description string  // A short description (optional)
    Path        string  // Current path of the file
    DeployPath  string  // Where we need to put this file (From the dotty.map)
}
