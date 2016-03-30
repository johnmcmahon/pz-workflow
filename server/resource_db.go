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

package server

import "github.com/venicegeo/pz-gocommon/elasticsearch"

type ResourceDB struct {
	server *Server
	Es     *elasticsearch.Client
	Esi    *elasticsearch.Index
}

func NewResourceDB(server *Server, es *elasticsearch.Client, esi *elasticsearch.Index) (*ResourceDB, error) {
	db := &ResourceDB{
		server: server,
		Es:     es,
		Esi:    esi,
	}

	_ = esi.Delete()

	err := esi.Create()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *ResourceDB) Flush() error {

	err := db.Esi.Flush()
	if err != nil {
		return LoggedError("ResourceDB.Flush failed: %s", err)
	}

	return nil
}
