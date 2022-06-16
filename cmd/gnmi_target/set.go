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
	// "encoding/json"
	// "fmt"
	// "io/fs"
	// "io/ioutil"
	// "strconv"

	"github.com/google/gnxi/utils/credentials"
	"github.com/openconfig/gnmi/proto/gnmi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	// dataStore "github.com/onosproject/storage-service/pkg/data_store"
	// types "github.com/onosproject/storage-service/pkg/types"
)

// Set overrides the Set func of gnmi.Target to provide user auth.
func (s *server) Set(ctx context.Context, req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		log.Infof("denied a Set request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}
	log.Info("allowed a Set request: %v", msg)

	updateStatus := "Failed"

	for _, update := range req.Update {
		if update.Path.Elem[0].Name == "Action" {
			if update.Path.Elem[0].Key["Action"] == "StoreGetReq" {
				log.Info("Storing get request now...")
				updateStatus = "Success"
			}
		}
	}

	updateResult := []*gnmi.UpdateResult{
		{
			Path: &gnmi.Path{
				Elem: []*gnmi.PathElem{
					{
						Name: "ActionResult",
						Key: map[string]string{
							"ActionResult": updateStatus,
						},
					},
				},
			},
		},
	}

	resp := &gnmi.SetResponse{
		Response: updateResult,
	}

	return resp, nil
}

// OLD CODE IN "Set"
// TODO: Refactor
// for _, update := range req.Update {
// 	if update.Path.Elem[0].Name == "Action" {
// 		if update.Path.Elem[0].Key["Action"] == "Change config" {
// 			updateResult = append(updateResult, s.updateConfig(update))
// 		} else if update.Path.Elem[0].Key["Action"] == "Store namespaces" {
// 			// TODO: store namespaces
// 			// extractNamespaces(update.Val.GetBytesVal())
//
// 			var schema Schema
// 			json.Unmarshal(update.Val.GetBytesVal(), &schema)
// 			schemaTree := getTreeStructure(schema)
//
// 			s.updateSchemaTreeList(schemaTree)
//
// 			s.storeSchemaTree(update.Val.GetBytesVal())
//
// 			updateResult = append(updateResult, &gnmi.UpdateResult{
// 				Path: &gnmi.Path{
// 					Elem: []*gnmi.PathElem{
// 						{
// 							Name: "Action",
// 							Key: map[string]string{
// 								"Action": "Successful",
// 							},
// 						},
// 					},
// 				},
// 			})
// 			// fmt.Println(path)
// 		}
// 	} else {
// 		fmt.Println("First element in path must be an action!")
// 	}
// }

// func (s *server) updateSchemaTreeList(schemaTree *SchemaTree) {
// 	alreadyStored := false
// 	for _, entry := range s.schemaTrees {
// 		if entry.Name == schemaTree.Name && entry.Namespace == schemaTree.Namespace {
// 			alreadyStored = true
// 		}
// 	}
// 	if !alreadyStored {
// 		s.schemaTrees = append(s.schemaTrees, schemaTree)
// 	} else {
// 		log.Info("SchemaTree already stored in list!")
// 	}
// }

// func (s *server) storeSchemaTree(schemaBytes []byte) {
// 	err := ioutil.WriteFile("configs_with_ns.json", schemaBytes, fs.ModeAppend.Perm())
// 	if err != nil {
// 		log.Warnf("Failed to write schemaTree to file! err: %v", err)
// 	} else {
// 		log.Info("Successfully wrote schemaTree to file!")
// 	}
// }

// // func extractNamespaces(bytes []byte) {
// // 	var schema Schema
// // 	json.Unmarshal(bytes, &schema)
//
// // 	// fmt.Println(schema)
// // 	schemaTree := getTreeStructure(schema)
//
// // 	// fmt.Println("#######################")
// // 	fmt.Println(schemaTree.Name)
// // 	fmt.Println("--------")
// // 	for _, child := range schemaTree.Children {
// // 		fmt.Print(" - ")
// // 		fmt.Print(child.Name)
// // 		fmt.Print(", ")
// // 		fmt.Println(child.Namespace)
// // 	}
// // }

type SchemaTree struct {
	Name      string
	Namespace string
	Children  []*SchemaTree
	Parent    *SchemaTree
	Value     string
}

// func getTreeStructure(schema Schema) *SchemaTree {
// 	var newTree *SchemaTree
// 	tree := &SchemaTree{}
// 	lastNode := ""
// 	for _, entry := range schema.Entries {
// 		if entry.Value == "" { // Directory
// 			if entry.Tag == "end" {
// 				if entry.Name != "data" {
// 					if lastNode != "leaf" {
// 						tree = tree.Parent
// 					}
// 					lastNode = ""
// 				}
// 			} else {
// 				newTree = &SchemaTree{Parent: tree}
// 				newTree.Name = entry.Name
// 				newTree.Namespace = entry.Namespace
// 				newTree.Parent.Children = append(newTree.Parent.Children, newTree)
//
// 				tree = newTree
// 			}
// 		} else { // Leaf
// 			newTree = &SchemaTree{Parent: tree}
// 			newTree.Name = entry.Name
// 			newTree.Value = entry.Value
// 			newTree.Parent.Children = append(newTree.Parent.Children, newTree)
//
// 			lastNode = "leaf"
// 		}
// 	}
// 	return tree
// }

// // func (s *server) updateConfig(update *gnmi.Update) *gnmi.UpdateResult {
// // 	config := dataStore.GetFullConfig()
//
// // 	var configObject types.ConfigObject
//
// // 	infoExists := false
// // 	for _, elem := range update.Path.Elem[1:] {
// // 		if elem.Name == "Info" {
// // 			infoExists = true
// // 			for _, confObj := range config.Devices {
// // 				if elem.Key["DeviceIp"] == confObj.DeviceIP {
// // 					configObject = confObj
// // 				}
// // 			}
// // 			break
// // 		}
// // 	}
// // 	if !infoExists {
// // 		fmt.Println("Could not update config as info element is missing!")
// // 		// TODO: Set gnmi.UpdateResult to be invalid
// // 		return nil
// // 	}
//
// // 	// TODO: Add mutex locks/semaphores on writing to the datastore
//
// // 	// TODO: update/create config for a given device
//
// // 	index := 0
//
// // 	// var config types.Config
// // 	// var conf types.ConfigObject
//
// // 	for _, elem := range update.Path.Elem[2:] {
// // 		switch elem.Name {
// // 		// case "Info":
// // 		// 	{
// // 		// 		conf.DeviceIP = elem.Key["DeviceIP"]
// // 		// 		conf.DeviceName = elem.Key["DeviceName"]
// // 		// 		conf.Protocol = elem.Key["Protocol"]
// // 		// 	}
// // 		case "Config" + strconv.Itoa(index):
// // 			{
// // 				// TODO: Get counters and build into Config object, then store that object
// // 				// as a config file, but also keep an up to date object in memory.
//
// // 				deviceCounters := s.getCounterData(elem, index)
//
// // 				// fmt.Println(deviceCounters)
//
// // 				err := s.modifyCounterData(&configObject, &deviceCounters)
// // 				if err != nil {
// // 					fmt.Println("Failed to modify counters!")
// // 				}
//
// // 				index++
// // 			}
// // 		default:
// // 			{
// // 				fmt.Println("Elem not recognized!")
// // 			}
// // 		}
// // 	}
//
// // 	return nil
// // }

// // func (s *server) modifyCounterData(confObj *types.ConfigObject, counters *[]types.DeviceCounter) error {
// // 	for _, oldConfObj := range confObj.Configs {
// // 		for _, oldCounter := range oldConfObj.Counter {
// // 			for _, newCounter := range *counters {
// // 				if oldCounter.Name == newCounter.Name {
// // 					// TODO: Replace old counter fields that new counter has, do not
// // 					// replace with empty values though.
//
// // 					break
// // 				}
// // 			}
// // 		}
// // 	}
//
// // 	return nil
// // }

// // func (s *server) getCounterData(elem *gnmi.PathElem, index int) []types.DeviceCounter {
// // 	indStr := strconv.Itoa(index)
//
// // 	var counters []types.DeviceCounter
//
// // 	var counter types.DeviceCounter
// // 	var err error
//
// // 	i := 0
// // 	for name, key := range elem.Key {
// // 		switch name {
// // 		case "Interval" + indStr:
// // 			{
// // 				// fmt.Println("Interval is: " + key)
// // 				counter.Interval, err = strconv.Atoi(key)
// // 				if err != nil {
// // 					fmt.Println("Failed to convert interval to int!")
// // 				}
// // 			}
// // 		case "Name" + indStr:
// // 			{
// // 				// fmt.Println("Name is: " + key)
// // 				counter.Name = key
// // 			}
// // 		case "Path" + indStr:
// // 			{
// // 				// fmt.Println("Path is: " + key)
// // 				counter.Path = key
// // 			}
// // 		default:
// // 			{
// // 				fmt.Println("Did not recognize the key!")
// // 			}
// // 		}
//
// // 		if i%3 == 2 {
// // 			fmt.Println(counter)
// // 			counters = append(counters, counter)
// // 		}
// // 		i++
// // 	}
//
// // 	return counters
// // }
