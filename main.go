// Copyright 2016, RadiantBlue Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"

	piazza "github.com/venicegeo/pz-gocommon"
	"github.com/venicegeo/pz-gocommon/elasticsearch"
	loggerPkg "github.com/venicegeo/pz-logger/lib"
	uuidgenPkg "github.com/venicegeo/pz-uuidgen/client"
	"github.com/venicegeo/pz-workflow/server"
)

func main() {

	required := []piazza.ServiceName{
		piazza.PzElasticSearch,
		piazza.PzLogger,
		piazza.PzGateway,
		piazza.PzUuidgen,
	}

	sys, err := piazza.NewSystemConfig(piazza.PzWorkflow, required)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := loggerPkg.NewClient(sys)
	if err != nil {
		log.Fatal(err)
	}

//	uuidgen, err := uuidgenPkg.NewMockUuidGenService(sys)
	uuidgen, err := uuidgenPkg.NewPzUuidGenService(sys)
	
	if err != nil {
		log.Fatal(err)
	}

	eventtypesIndex, err := elasticsearch.NewIndex(sys, "eventtypes")
	if err != nil {
		log.Fatal(err)
	}
	eventsIndex, err := elasticsearch.NewIndex(sys, "events")
	if err != nil {
		log.Fatal(err)
	}
	triggersIndex, err := elasticsearch.NewIndex(sys, "triggers")
	if err != nil {
		log.Fatal(err)
	}
	alertsIndex, err := elasticsearch.NewIndex(sys, "alerts")
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("pz-workflow starting...")

	// start server
	routes, err := server.CreateHandlers(sys, logger, uuidgen,
		eventtypesIndex, eventsIndex, triggersIndex, alertsIndex)
	if err != nil {
		log.Fatal(err)
	}

	// _ = sys.StartServer(routes)

	done := sys.StartServer(routes)

	err = <-done
	if err != nil {
		log.Fatal(err)
	}
}
