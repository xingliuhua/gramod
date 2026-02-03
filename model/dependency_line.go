package model

// Module represents a go module with its path and version.
type Module struct {
	Path    string
	Version string
}

// DependencyMap is an adjacency list representation of module dependencies.
// Key = parent module, value = slice of direct dependencies.
type DependencyMap map[Module][]Module
