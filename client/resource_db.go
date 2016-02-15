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

package client

import (
	"encoding/json"
	"github.com/venicegeo/pz-gocommon"
)

var resourceID = 1

func NewResourceID() Ident {
	id := NewIdentFromInt(resourceID)
	resourceID++
	return Ident("R" + string(id))
}

//type Resource interface {
//	GetId() Ident
//	SetId(Ident)
//}

type ResourceDB struct {
	es       *piazza.ElasticSearchService
	index    string
	typename string
}

func NewResourceDB(es *piazza.ElasticSearchService, index string, typename string) (*ResourceDB, error) {
	db := &ResourceDB{
		es:       es,
		index:    index,
		typename: typename,
	}

	err := es.DeleteIndex(index)
	if err != nil {
		return nil, err
	}

	err = es.CreateIndex(index)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *ResourceDB) PostData(obj interface{}, id Ident) (Ident, error) {

	_, err := db.es.PostData(db.index, db.typename, id.String(), obj)
	if err != nil {
		return NoIdent, err
	}

	err = db.es.FlushIndex(db.index)
	if err != nil {
		return NoIdent, err
	}

	return id, nil
}

func (db *ResourceDB) GetAll() ([]*json.RawMessage, error) {
	searchResult, err := db.es.SearchByMatchAll(db.index)
	if err != nil {
		return nil, err
	}

	if searchResult.Hits == nil {
		return nil, nil
	}

	raws := make([]*json.RawMessage,searchResult.TotalHits())

	for i, hit := range searchResult.Hits.Hits {
		raws[i] = hit.Source
	}

	return raws, nil
}

func (db *ResourceDB) GetById(id Ident, obj interface{}) (bool, error) {

	getResult, err := db.es.GetById(db.index, id.String())
	if err != nil {
		return false, err
	}

	if !getResult.Found {
		return false, nil
	}

	src := getResult.Source
	err = json.Unmarshal(*src, obj)
	if err != nil {
		return true, err
	}
	return true, nil
}

func (db *ResourceDB) DeleteByID(id string) (bool, error) {
	res, err := db.es.DeleteById(db.index, db.typename, id)
	if err != nil {
		return false, err
	}

	err = db.es.FlushIndex(db.index)
	if err != nil {
		return false, err
	}

	return res.Found, nil
}