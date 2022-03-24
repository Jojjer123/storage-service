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
	"time"

	"github.com/google/gnxi/utils/credentials"
	"github.com/openconfig/gnmi/proto/gnmi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	// "k8s.io/apimachinery/pkg/util/json"
)

// Get overrides the Get func of gnmi.Target to provide user auth.
func (s *server) Get(ctx context.Context, req *gnmi.GetRequest) (*gnmi.GetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		log.Infof("denied a Get request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}

	log.Infof("allowed a Get request: %v", msg)

	notifications := make([]*gnmi.Notification, 1)
	prefix := req.GetPrefix()
	ts := time.Now().UnixNano()

	// COMMENTED FOR TESTING SCHEMA
	update := make([]*gnmi.Update, 1)

	// update[0] = &gnmi.Update{
	// 	Value: &gnmi.Value{
	// 		Value: dataStore.GetConfig(req),
	// 	},
	// }

	// TEMPORARY CHANGES
	schema := netconfConv(xmlString)

	// jsonDump, err := json.Marshal(schema)
	// if err != nil {
	// 	fmt.Println("Failed to marshal schema")
	// 	fmt.Println(err)
	// }

	// gob.Register(map[string]interface{}{})

	// byteVal := []byte{schema}

	// // gob.Register(map[string]interface{}{})
	// // gob.Register(Schema{})
	// // gob.Register(xml.CharData{})

	// // byteVal := bytes.Buffer{}
	// // enc := gob.NewEncoder(&byteVal)
	// // err := enc.Encode(schema)
	// // if err != nil {
	// // 	fmt.Println("Failed to encode schema!")
	// // 	fmt.Println(err)
	// // }

	jsonDump, err := json.Marshal(schema)
	if err != nil {
		fmt.Println("Failed to serialize schema!")
		fmt.Println(err)
	}

	update[0] = &gnmi.Update{Val: &gnmi.TypedValue{Value: &gnmi.TypedValue_BytesVal{BytesVal: jsonDump}}} // byteVal.Bytes()}}}

	// END OF TEMPORARY CHANGES

	notifications[0] = &gnmi.Notification{
		Timestamp: ts,
		Prefix:    prefix,
		Update:    update,
	}

	resp := &gnmi.GetResponse{Notification: notifications}

	return resp, nil
	// return s.Server.Get(ctx, req)
}

// func (s Schema) MarshalBinary() (_ []byte, err error) {
// 	var buf bytes.Buffer
// 	enc := gob.NewEncoder(&buf)
// 	enc.Encode(s.Entry)
// 	enc.Encode(s.Children)

// 	if s.Parent.Entry == nil {
// 		return buf.Bytes(), nil
// 	}

// 	// isCyclic := s.Parent.Entry != nil && s.Parent.Parent == &s
// 	// enc.Encode(isCyclic)
// 	// if isCyclic {
// 	// 	s.Parent.Parent = nil
// 	// 	err = enc.Encode(s.Parent)
// 	// 	p.Q.P = p
// 	// } else {
// 	// 	err = enc.Encode(p.Q)
// 	// }

// 	return buf.Bytes(), err
// }

// func encodeObject(object interface{}) []byte {
// 	buf := bytes.Buffer{}

// 	enc := gob.NewEncoder(&buf)
// 	err := enc.Encode(object)
// 	if err != nil {
// 		fmt.Println("failed to encode object!")
// 		fmt.Println(err)
// 	}

// 	return buf.Bytes()
// }
