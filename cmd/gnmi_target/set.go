// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/gnxi/utils/credentials"
	"github.com/openconfig/gnmi/proto/gnmi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	dataStore "github.com/onosproject/storage-service/pkg/data_store"
	types "github.com/onosproject/storage-service/pkg/types"
)

// Set overrides the Set func of gnmi.Target to provide user auth.
func (s *server) Set(ctx context.Context, req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		log.Infof("denied a Set request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}
	fmt.Println("allowed a Set request")

	var updateResult []*gnmi.UpdateResult

	for _, update := range req.Update {
		if update.Path.Elem[0].Name == "Action" {
			if update.Path.Elem[0].Key["Action"] == "Change config" {
				updateResult = append(updateResult, s.updateConfig(update))
			} else if update.Path.Elem[0].Key["Action"] == "Store namespaces" {
				// TODO: store namespaces
				// extractNamespaces(update.Val.GetBytesVal())

				var schema Schema
				json.Unmarshal(update.Val.GetBytesVal(), &schema)
				schemaTree := getTreeStructure(schema)

				pathElements := getNamespacesForPath([]*gnmi.PathElem{
					{
						Name: "interfaces",
						Key:  map[string]string{},
					},
					{
						Name: "interface",
						Key:  map[string]string{},
					},
				}, schemaTree.Children)

				fmt.Println(pathElements)
			}
		} else {
			fmt.Println("First element in path must be an action!")
		}
	}

	resp := &gnmi.SetResponse{
		Response:  updateResult,
		Timestamp: time.Now().UnixNano(),
	}

	return resp, nil
}

// TODO: Add key reading in path so that specific elements based on keys can be used.
func getNamespacesForPath(pathElems []*gnmi.PathElem, schemaTreeChildren []*SchemaTree) []*gnmi.PathElem {
	var newPathElems []*gnmi.PathElem
	childFound := false
	if len(pathElems) > 0 {
		for _, child := range schemaTreeChildren {
			// fmt.Println(child.Name)
			if pathElems[0].Name == child.Name {
				if child.Namespace != "" {
					pathElems[0].Key["namespace"] = child.Namespace
					// fmt.Printf("Added namespace \"%s\" for child: %s\n", child.Namespace, child.Name)
				}
				// else {
				// 	fmt.Printf("No namespace for elem %s\n", child.Name)
				// }
				childFound = true
				newPathElems = append(newPathElems, pathElems[0])
				return append(newPathElems, getNamespacesForPath(pathElems[1:], child.Children)[0])
				// break
			}
		}

		if !childFound {
			fmt.Printf("Could not find path element: %s", pathElems[0].Name)
			return nil
		}
	}

	// fmt.Println("No more elements!")

	return nil
}

// func extractNamespaces(bytes []byte) {
// 	var schema Schema
// 	json.Unmarshal(bytes, &schema)

// 	// fmt.Println(schema)
// 	schemaTree := getTreeStructure(schema)

// 	// fmt.Println("#######################")
// 	fmt.Println(schemaTree.Name)
// 	fmt.Println("--------")
// 	for _, child := range schemaTree.Children {
// 		fmt.Print(" - ")
// 		fmt.Print(child.Name)
// 		fmt.Print(", ")
// 		fmt.Println(child.Namespace)
// 	}
// }

type SchemaTree struct {
	Name      string
	Namespace string
	Children  []*SchemaTree
	Parent    *SchemaTree
	Value     string
}

// TODO: add pointer that traverse the tree based on tags, use that pointer to
// get correct parents.
func getTreeStructure(schema Schema) *SchemaTree {
	var newTree *SchemaTree
	tree := &SchemaTree{}
	lastNode := ""
	for _, entry := range schema.Entries {
		// fmt.Println("-------------------")
		// if index == 0 {
		// newTree = &SchemaTree{Parent: tree}
		// newTree.Name = entry.Name
		// newTree.Namespace = entry.Namespace
		// fmt.Println(tree.Name)
		// tree = &SchemaTree{Parent: tree}
		// continue
		// }
		if entry.Value == "" { // Directory
			if entry.Tag == "end" {
				if entry.Name != "data" {
					if lastNode != "leaf" {
						// fmt.Println(tree.Name)
						tree = tree.Parent
					}
					lastNode = ""
					// continue
				}
			} else {

				newTree = &SchemaTree{Parent: tree}

				newTree.Name = entry.Name
				newTree.Namespace = entry.Namespace
				newTree.Parent.Children = append(newTree.Parent.Children, newTree)

				tree = newTree

				// fmt.Print(tree.Name)
				// fmt.Print(", ")
				// fmt.Println(tree.Parent.Name)

				// tree = &SchemaTree{Parent: tree}
			}
		} else { // Leaf
			newTree = &SchemaTree{Parent: tree}

			newTree.Name = entry.Name
			newTree.Value = entry.Value
			newTree.Parent.Children = append(newTree.Parent.Children, newTree)

			// fmt.Print(newTree.Name)
			// fmt.Print(", ")
			// fmt.Println(newTree.Parent.Name)
			// fmt.Println(newTree.Value)

			lastNode = "leaf"
		}
		// fmt.Println("-------------------")
		// fmt.Print("name: ")
		// fmt.Print(tree.Name)
		// if tree.Name != "data" {
		// 	fmt.Print(", parent: ")
		// 	fmt.Println(tree.Parent.Name)

		// 	// fmt.Println("#######")
		// 	// for i, child := range tree.Parent.Children {
		// 	// 	if i < 10 {
		// 	// 		fmt.Print(child.Name)
		// 	// 		fmt.Print(", ")
		// 	// 	}
		// 	// }
		// 	// fmt.Println("\n******")
		// 	// for j, child := range tree.Children {
		// 	// 	if j < 10 {
		// 	// 		fmt.Print(child.Name)
		// 	// 		fmt.Print(": ")
		// 	// 		fmt.Print(child.Value)
		// 	// 		fmt.Print(", ")
		// 	// 	}
		// 	// }
		// } else {
		// 	// fmt.Println("")
		// 	// for _, child := range tree.Children {
		// 	// 	fmt.Print(child.Name)
		// 	// 	fmt.Print(" | ")
		// 	// }
		// }
		// fmt.Println("")
		// fmt.Println(entry)
		// fmt.Println(tree.Namespace)
		// fmt.Println("###################")
	}
	return tree
}

func (s *server) updateConfig(update *gnmi.Update) *gnmi.UpdateResult {
	config := dataStore.GetFullConfig()

	var configObject types.ConfigObject

	infoExists := false
	for _, elem := range update.Path.Elem[1:] {
		if elem.Name == "Info" {
			infoExists = true
			for _, confObj := range config.Devices {
				if elem.Key["DeviceIp"] == confObj.DeviceIP {
					configObject = confObj
				}
			}
			break
		}
	}
	if !infoExists {
		fmt.Println("Could not update config as info element is missing!")
		// TODO: Set gnmi.UpdateResult to be invalid
		return nil
	}

	// TODO: Add mutex locks/semaphores on writing to the datastore

	// TODO: update/create config for a given device

	index := 0

	// var config types.Config
	// var conf types.ConfigObject

	for _, elem := range update.Path.Elem[2:] {
		switch elem.Name {
		// case "Info":
		// 	{
		// 		conf.DeviceIP = elem.Key["DeviceIP"]
		// 		conf.DeviceName = elem.Key["DeviceName"]
		// 		conf.Protocol = elem.Key["Protocol"]
		// 	}
		case "Config" + strconv.Itoa(index):
			{
				// TODO: Get counters and build into Config object, then store that object
				// as a config file, but also keep an up to date object in memory.

				deviceCounters := s.getCounterData(elem, index)

				// fmt.Println(deviceCounters)

				err := s.modifyCounterData(&configObject, &deviceCounters)
				if err != nil {
					fmt.Println("Failed to modify counters!")
				}

				index++
			}
		default:
			{
				fmt.Println("Elem not recognized!")
			}
		}
	}

	return nil
}

func (s *server) modifyCounterData(confObj *types.ConfigObject, counters *[]types.DeviceCounters) error {
	for _, oldConfObj := range confObj.Configs {
		for _, oldCounter := range oldConfObj.Counter {
			for _, newCounter := range *counters {
				if oldCounter.Name == newCounter.Name {
					// TODO: Replace old counter fields that new counter has, do not
					// replace with empty values though.

					break
				}
			}
		}
	}

	return nil
}

func (s *server) getCounterData(elem *gnmi.PathElem, index int) []types.DeviceCounters {
	indStr := strconv.Itoa(index)

	var counters []types.DeviceCounters

	var counter types.DeviceCounters
	var err error

	i := 0
	for name, key := range elem.Key {
		switch name {
		case "Interval" + indStr:
			{
				// fmt.Println("Interval is: " + key)
				counter.Interval, err = strconv.Atoi(key)
				if err != nil {
					fmt.Println("Failed to convert interval to int!")
				}
			}
		case "Name" + indStr:
			{
				// fmt.Println("Name is: " + key)
				counter.Name = key
			}
		case "Path" + indStr:
			{
				// fmt.Println("Path is: " + key)
				counter.Path = key
			}
		default:
			{
				fmt.Println("Did not recognize the key!")
			}
		}

		if i%3 == 2 {
			fmt.Println(counter)
			counters = append(counters, counter)
		}
		i++
	}

	return counters
}
