package slinga

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mattn/go-zglob"
	"sort"
	. "github.com/Frostman/aptomi/pkg/slinga/maputil"
	. "github.com/Frostman/aptomi/pkg/slinga/log"
	. "github.com/Frostman/aptomi/pkg/slinga/fileio"
)

/*
	This file declares all the necessary structures for Dependencies (User "wants" Service)
*/

// Dependency in a form <UserID> requested <Service> (and provided additional <Labels>)
type Dependency struct {
	Enabled bool
	ID      string
	UserID  string
	Service string
	Labels  map[string]string

	// This fields are populated when dependency gets resolved
	Resolved   bool
	ServiceKey string
}

// UnmarshalYAML is a custom unmarshaller for Dependency, which sets Enabled to True before unmarshalling
func (s *Dependency) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type Alias Dependency
	instance := Alias{Enabled: true}
	if err := unmarshal(&instance); err != nil {
		return err
	}
	*s = Dependency(instance)
	return nil
}

// GlobalDependencies represents the list of global dependencies (see the definition above)
// TODO: during serialization there is data duplication (as both fields get serialized). should prob avoid this
type GlobalDependencies struct {
	// dependencies <service> -> list of dependencies
	DependenciesByService map[string][]*Dependency

	// dependencies <id> -> dependency
	DependenciesByID map[string]*Dependency
}

func (src *GlobalDependencies) count() int {
	return CountElements(src.DependenciesByID)
}

// NewGlobalDependencies creates and initializes a new empty list of global dependencies
func NewGlobalDependencies() GlobalDependencies {
	return GlobalDependencies{
		DependenciesByService: make(map[string][]*Dependency),
		DependenciesByID:      make(map[string]*Dependency),
	}
}

// Apply set of transformations to labels
func (dependency *Dependency) getLabelSet() LabelSet {
	return LabelSet{Labels: dependency.Labels}
}

// Append a single dependency to an existing object
func (src GlobalDependencies) appendDependency(dependency *Dependency) {
	if len(dependency.ID) <= 0 {
		Debug.WithFields(log.Fields{
			"dependency": dependency,
		}).Panic("Empty dependency ID")
	}
	src.DependenciesByService[dependency.Service] = append(src.DependenciesByService[dependency.Service], dependency)
	src.DependenciesByID[dependency.ID] = dependency
}

// Copy the whole structure with dependencies
func (src GlobalDependencies) makeCopy() GlobalDependencies {
	result := NewGlobalDependencies()
	for _, v := range src.DependenciesByID {
		result.appendDependency(v)
	}
	return result
}

// LoadDependenciesFromDir loads all dependencies from a given directory
func LoadDependenciesFromDir(baseDir string) GlobalDependencies {
	// read all services
	files, _ := zglob.Glob(GetAptomiObjectFilePatternYaml(baseDir, TypeDependencies))
	sort.Strings(files)
	result := NewGlobalDependencies()
	for _, fileName := range files {
		t := loadDependenciesFromFile(fileName)
		for _, d := range t {
			if d.Enabled {
				result.appendDependency(d)
			}
		}
	}
	return result
}
