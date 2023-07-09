package domain

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	Development  = "dev"
	Homologation = "hmg"
	Production   = "prd"
)

type Version struct {
	Enviroment Enviroment `yaml:"version"`
}

type Enviroment struct {
	Dev string `yaml:"dev"`
	Hmg string `yaml:"hmg"`
	Prd string `yaml:"prd"`
}

func (v Version) getEnviroment(enviroment string) string {
	switch enviroment {
	case Development:
		return v.Enviroment.Dev
	case Homologation:
		return v.Enviroment.Hmg
	case Production:
		return v.Enviroment.Prd
	default:
		panic("Invalid environment type. Options: dev, hmg, prd")
	}
}

func (v *Version) setEnviroment(enviroment, version string) {
	switch enviroment {
	case Development:
		v.Enviroment.Dev = version
	case Homologation:
		v.Enviroment.Hmg = version
	case Production:
		v.Enviroment.Prd = version
	default:
		panic("Invalid environment type. Options: dev, hmg, prd")
	}
}

func (v *Version) IncrementByEnviroment(index int, enviroment string) string {
	parts := strings.Split(v.getEnviroment(enviroment), ".")
	if len(parts) < 3 {
		panic("Invalid version in format 'x.y.z'")
	}

	newValue, err := strconv.Atoi(parts[index])
	if err != nil {
		panic("Reported version must be a valid integer")
	}
	newValue++

	return v.format(index, newValue, parts, enviroment)
}

func (v *Version) format(index, newValue int, parts []string, enviroment string) string {
	major := 0
	minor := 1

	var version string
	if index == major {
		version = fmt.Sprintf("%d.0.0", newValue)
	} else if index == minor {
		version = fmt.Sprintf("%s.%d.0", parts[0], newValue)
	} else {
		version = fmt.Sprintf("%s.%s.%d", parts[0], parts[1], newValue)
	}
	v.setEnviroment(enviroment, version)
	return version
}

func (v Version) ToByte() []byte {
	newYaml, err := yaml.Marshal(&v)
	if err != nil {
		panic("Error encoding YML file: " + err.Error())
	}
	return newYaml
}

func (v Version) Extract(file *os.File) *Version {
	content, err := io.ReadAll(file)
	if err != nil {
		panic("Error reading file: " + err.Error())
	}

	err = yaml.Unmarshal(content, &v)
	if err != nil {
		panic("Error parsing YML file: " + err.Error())
	}

	fmt.Printf("\nVersion in Development: %s", v.Enviroment.Dev)
	fmt.Printf("\nVersion in Homologation: %s", v.Enviroment.Hmg)
	fmt.Printf("\nVersion in Production: %s", v.Enviroment.Prd)
	fmt.Println("\n\n---Done---")

	return &v
}
