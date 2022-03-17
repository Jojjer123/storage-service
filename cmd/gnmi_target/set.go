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
	"fmt"
	"strconv"
	"time"

	"github.com/google/gnxi/utils/credentials"
	"github.com/openconfig/gnmi/proto/gnmi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	// fmt.Println(req.Update[0].Path)

	var updateResult []*gnmi.UpdateResult

	for _, update := range req.Update {
		if update.Path.Elem[0].Name == "Action" {
			if update.Path.Elem[0].Key["Action"] == "Change config" {
				updateResult = append(updateResult, s.updateConfig(update))
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

func (s *server) updateConfig(update *gnmi.Update) *gnmi.UpdateResult {
	index := 0
	for _, elem := range update.Path.Elem[1:] {
		switch elem.Name {
		case "Info":
			{

			}
		case "Config" + strconv.Itoa(index):
			{
				// TODO: Get counters and build into Config object, then store that object
				// as a config file, but also keep an up to date object in memory.

				deviceCounters := s.getCounterData(elem)

				fmt.Println(deviceCounters)

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

func (s *server) getCounterData(*gnmi.PathElem) []types.DeviceCounters {

	return []types.DeviceCounters{}
}
