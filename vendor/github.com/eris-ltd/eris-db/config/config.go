// Copyright 2017 Monax Industries Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"bytes"
	"fmt"
	"text/template"
)

type ConfigServiceGeneral struct {
	ChainImageName      string
	UseDataContainer    bool
	ExportedPorts       string
	ContainerEntrypoint string
}

// TODO: [ben] increase the configurability upon need
type ConfigChainGeneral struct {
	AssertChainId       string
	ErisdbMajorVersion  uint8
	ErisdbMinorVersion  uint8
	GenesisRelativePath string
}

type ConfigChainModule struct {
	Name               string
	MajorVersion       uint8
	MinorVersion       uint8
	ModuleRelativeRoot string
}

type ConfigTendermint struct {
	Moniker  string
	Seeds    string
	FastSync bool
}

var serviceGeneralTemplate *template.Template
var chainGeneralTemplate *template.Template
var chainConsensusTemplate *template.Template
var chainApplicationManagerTemplate *template.Template
var tendermintTemplate *template.Template

func init() {
	var err error
	if serviceGeneralTemplate, err = template.New("serviceGeneral").Parse(sectionServiceGeneral); err != nil {
		panic(err)
	}
	if chainGeneralTemplate, err = template.New("chainGeneral").Parse(sectionChainGeneral); err != nil {
		panic(err)
	}
	if chainConsensusTemplate, err = template.New("chainConsensus").Parse(sectionChainConsensus); err != nil {
		panic(err)
	}
	if chainApplicationManagerTemplate, err = template.New("chainApplicationManager").Parse(sectionChainApplicationManager); err != nil {
		panic(err)
	}
	if tendermintTemplate, err = template.New("tendermint").Parse(sectionTendermint); err != nil {
		panic(err)
	}
}

// NOTE: [ben] for 0.12.0-rc3 we only have a single configuration path
// with Tendermint in-process as the consensus engine and ErisMint
// in-process as the application manager, so we hard-code the few
// parameters that are already templated.
// Let's learn to walk before we can run.
func GetConfigurationFileBytes(chainId, moniker, seeds string, chainImageName string,
	useDataContainer bool, exportedPortsString, containerEntrypoint string) ([]byte, error) {

	erisdbService := &ConfigServiceGeneral{
		ChainImageName:      chainImageName,
		UseDataContainer:    useDataContainer,
		ExportedPorts:       exportedPortsString,
		ContainerEntrypoint: containerEntrypoint,
	}
	erisdbChain := &ConfigChainGeneral{
		AssertChainId:       chainId,
		ErisdbMajorVersion:  uint8(0),
		ErisdbMinorVersion:  uint8(16),
		GenesisRelativePath: "genesis.json",
	}
	chainConsensusModule := &ConfigChainModule{
		Name:               "tendermint",
		MajorVersion:       uint8(0),
		MinorVersion:       uint8(8),
		ModuleRelativeRoot: "tendermint",
	}
	chainApplicationManagerModule := &ConfigChainModule{
		Name:               "erismint",
		MajorVersion:       uint8(0),
		MinorVersion:       uint8(16),
		ModuleRelativeRoot: "erismint",
	}
	tendermintModule := &ConfigTendermint{
		Moniker:  moniker,
		Seeds:    seeds,
		FastSync: false,
	}

	// NOTE: [ben] according to StackOverflow appending strings with copy is
	// more efficient than bytes.WriteString, but for readability and because
	// this is not performance critical code we opt for bytes, which is
	// still more efficient than + concatentation operator.
	var buffer bytes.Buffer

	// write copyright header
	buffer.WriteString(headerCopyright)

	// write section [service]
	if err := serviceGeneralTemplate.Execute(&buffer, erisdbService); err != nil {
		return nil, fmt.Errorf("Failed to write template service general for %s: %s",
			chainId, err)
	}
	// write section for service dependencies; this is currently a static section
	// with a fixed dependency on eris-keys
	buffer.WriteString(sectionServiceDependencies)

	// write section [chain]
	if err := chainGeneralTemplate.Execute(&buffer, erisdbChain); err != nil {
		return nil, fmt.Errorf("Failed to write template chain general for %s: %s",
			chainId, err)
	}

	// write separator chain consensus
	buffer.WriteString(separatorChainConsensus)
	// write section [chain.consensus]
	if err := chainConsensusTemplate.Execute(&buffer, chainConsensusModule); err != nil {
		return nil, fmt.Errorf("Failed to write template chain consensus for %s: %s",
			chainId, err)
	}

	// write separator chain application manager
	buffer.WriteString(separatorChainApplicationManager)
	// write section [chain.consensus]
	if err := chainApplicationManagerTemplate.Execute(&buffer,
		chainApplicationManagerModule); err != nil {
		return nil, fmt.Errorf("Failed to write template chain application manager for %s: %s",
			chainId, err)
	}

	// write separator servers
	buffer.WriteString(separatorServerConfiguration)
	// TODO: [ben] upon necessity replace this with template too
	// write static section servers
	buffer.WriteString(sectionServers)

	// write separator modules
	buffer.WriteString(separatorModules)

	// write section module Tendermint
	if err := tendermintTemplate.Execute(&buffer, tendermintModule); err != nil {
		return nil, fmt.Errorf("Failed to write template tendermint for %s, moniker %s: %s",
			chainId, moniker, err)
	}

	// write static section erismint
	buffer.WriteString(sectionErisMint)

	return buffer.Bytes(), nil
}

func GetExampleConfigFileBytes() ([]byte, error) {
	return GetConfigurationFileBytes(
		"simplechain",
		"delectable_marmot",
		"192.168.168.255",
		"db:latest",
		true,
		"46657",
		"eris-db")
}
