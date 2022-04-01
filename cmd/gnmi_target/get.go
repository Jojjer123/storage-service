package main

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/google/gnxi/utils/credentials"
	"github.com/openconfig/gnmi/proto/gnmi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get overrides the Get func of gnmi.Target to provide user auth.
func (s *server) Get(ctx context.Context, req *gnmi.GetRequest) (*gnmi.GetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		log.Infof("denied a Get request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}

	log.Infof("allowed a Get request: %v", msg)

	notifications := []*gnmi.Notification{
		{
			Update: []*gnmi.Update{},
		},
	}

	switch req.Type {
	case 4:
		for _, path := range req.Path {
			updatedPath := &gnmi.Path{}
			getNamespacesForPath(updatedPath,
				path.Elem, s.schemaTrees[0].Children)

			notifications[0].Update = append(notifications[0].Update, &gnmi.Update{Path: updatedPath})
		}
	default:
		log.Warn("Did not recognize requested type!")
	}

	printSchemaTree(s.schemaTrees[0])

	// COMMENTED FOR TESTING SCHEMA

	// update[0] = &gnmi.Update{
	// 	Value: &gnmi.Value{
	// 		Value: dataStore.GetConfig(req),
	// 	},
	// }

	// TEMPORARY CHANGES
	// schema := netconfConv(xmlString)

	// jsonDump, err := json.Marshal(schema)
	// if err != nil {
	// 	fmt.Println("Failed to serialize schema!")
	// 	fmt.Println(err)
	// }

	// update[0] = &gnmi.Update{Val: &gnmi.TypedValue{Value: &gnmi.TypedValue_BytesVal{BytesVal: jsonDump}}}

	// END OF TEMPORARY CHANGES

	resp := &gnmi.GetResponse{Notification: notifications}

	return resp, nil
}

func printSchemaTree(schemaTree *SchemaTree) {
	log.Infof("%s - %s - %v", schemaTree.Parent.Name, schemaTree.Name, schemaTree.Namespace)
	for _, child := range schemaTree.Children {
		printSchemaTree(child)
	}
}

// TODO: Add key reading in path so that specific elements based on keys can be used.
func getNamespacesForPath(path *gnmi.Path, pathElems []*gnmi.PathElem, schemaTreeChildren []*SchemaTree) {
	childFound := false

	if len(pathElems) <= 0 {
		return
	}

	for _, child := range schemaTreeChildren {
		// log.Infof("Current child is: %v", child)
		if len(pathElems[0].Key) > 0 {
			var key, val string
			for key, val = range pathElems[0].Key {
				if key != "namespace" {
					break
				}
			}
			// log.Infof("Found key: %s and value: %s", key, val)
			if pathElems[0].Name == child.Name {
				keyChildFound := false
				var keyChild *SchemaTree
				for _, keyChild = range child.Children {
					if val == keyChild.Value {
						// log.Infof("Found key child with name: %s and value: %s", keyChild.Name, keyChild.Value)
						keyChildFound = true
						break
					}
				}

				if keyChildFound {
					// log.Info("Adding namespace if there are one for the current child...")
					if child.Namespace != "" {
						// log.Infof("Child %v has namespace %s", child, child.Namespace)
						if pathElems[0].Key == nil {
							pathElems[0].Key = map[string]string{}
						}
						pathElems[0].Key["namespace"] = child.Namespace
					}
					childFound = true
					path.Elem = append(path.Elem, pathElems[0])
					getNamespacesForPath(path, pathElems[1:], child.Children)
					break
				} else {
					// log.Infof("%s was not correct child", child.Name)
				}
			}
		} else {
			if pathElems[0].Name == child.Name {
				// log.Infof("Child with name %s does not have any keys", child.Name)
				if child.Namespace != "" {
					if pathElems[0].Key == nil {
						pathElems[0].Key = map[string]string{}
					}
					pathElems[0].Key["namespace"] = child.Namespace
				}
				childFound = true
				path.Elem = append(path.Elem, pathElems[0])
				getNamespacesForPath(path, pathElems[1:], child.Children)
				break
			}
		}
	}

	if !childFound {
		fmt.Printf("Could not find path element: %s", pathElems[0].Name)
	}
}

// Converts XML to a Schema containing a slice of all the elements with namespaces and values.
// TODO: Add "searching" to filter out all data except what the path is requesting.
func netconfConv(xmlString string) *Schema {
	decoder := xml.NewDecoder(strings.NewReader(xmlString))
	schema := &Schema{}

	var newEntry *SchemaEntry
	var lastNamespace string

	index := 0
	for {
		tok, _ := decoder.Token()

		if tok == nil {
			return schema
		}

		switch tokType := tok.(type) {
		case xml.StartElement:
			newEntry = &SchemaEntry{}
			newEntry.Name = tokType.Name.Local

			if index > 0 {
				if tokType.Name.Space != lastNamespace {
					lastNamespace = tokType.Name.Space
					newEntry.Namespace = lastNamespace
				}
				newEntry.Tag = "start"
			} else {
				lastNamespace = tokType.Name.Space
				newEntry.Namespace = lastNamespace
			}

			schema.Entries = append(schema.Entries, *newEntry)
			index++

		case xml.EndElement:
			newEntry = &SchemaEntry{}
			newEntry.Name = tokType.Name.Local
			newEntry.Tag = "end"
			schema.Entries = append(schema.Entries, *newEntry)
			index++

		case xml.CharData:
			bytes := xml.CharData(tokType)
			schema.Entries[index-1].Value = string([]byte(bytes))

		default:
			fmt.Print(", was not recognized with type: ")
			fmt.Println(tokType)
		}
	}
}

// type Schema struct {
// 	Entry    *SchemaEntry // yang.Entry
// 	Children map[string]interface{}
// 	Parent   *Schema
// }

type Schema struct {
	Entries []SchemaEntry
}

type SchemaEntry struct {
	Name      string
	Tag       string
	Namespace string
	Value     string
}
