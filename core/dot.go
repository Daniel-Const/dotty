package core 


// Represents a Dot file
type Dot struct {
    Path        string  // Current path of the file
    DeployPath  string  // Where we need to put this file (From the dotty.map)
    IsDir       bool
}
