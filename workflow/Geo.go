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

package workflow

import (
	"encoding/json"
	"regexp"

	"github.com/venicegeo/pz-gocommon/gocommon"
)

const typeGeometryCollection = "geometrycollection"
const typePoint = "point"
const typeLineString = "linestring"
const typePolygon = "polygon"
const typeMultiPoint = "multipoint"
const typeMultiLineString = "multilinestring"
const typeMultiPolygon = "multipolygon"
const typeEnvelope = "envelope"
const typeCircle = "circle"

const orientationRegex = `^((right)|(ccw)|(counterclockwise)|(left)|(cw)|(clockwise))$`
const precisionRegex = `^((in)|(inch)|(yd)|(yard)|(mi)|(miles)|(km)|(kilometers)|(m)|(meters)|(cm)|(centimeters)|(mm)|(millimeters))$`
const distanceRegex = `^(([1-9][0-9]*)((in)|(inch)|(yd)|(yard)|(mi)|(miles)|(km)|(kilometers)|(m)|(meters)|(cm)|(centimeters)|(mm)|(millimeters)|$))$`

type Geo_Point struct {
	Lon float64 `json:"lon" binding:"required"`
	Lat float64 `json:"lat" binding:"required"`
}

func (p *Geo_Point) valid() bool { //TODO
	return true
}

type Geo_Shape struct {
	Type             interface{} `json:"type"`               //string
	Coordinates      interface{} `json:"coordinates"`        //interface{}
	Geometries       interface{} `json:"geometries"`         //interface{}
	Tree             interface{} `json:"tree"`               //string
	Precision        interface{} `json:"precision"`          //string
	TreeLevels       interface{} `json:"tree_levels"`        //string
	Strategy         interface{} `json:"strategy"`           //string
	DistanceErrorPct interface{} `json:"distance_error_pct"` //float64
	Orientation      interface{} `json:"orientation"`        //string
	PointsOnly       interface{} `json:"points_only"`        //bool
	Radius           interface{} `json:"radius"`             //string
}

type geo_GeometryCollection []Geo_Shape
type geo_Sub_Point []interface{}
type geo_LineString []geo_Sub_Point
type geo_Polygon [][]geo_Sub_Point
type geo_MultiPoint []geo_Sub_Point
type geo_MultiLineString []geo_LineString
type geo_MultiPolygon []geo_Polygon
type geo_Envelope []geo_Sub_Point
type geo_Circle geo_Sub_Point

func NewDefaultGeo_Shape() Geo_Shape {
	return Geo_Shape{Tree: "geohash", Precision: "meters", TreeLevels: "50m", Strategy: "recursive", DistanceErrorPct: 0.025, Orientation: "ccw", PointsOnly: false}
}

func (gs *Geo_Shape) valid() (bool, error) {
	if gs.Tree != nil {
		if ok, err := gs.validTree(gs.Tree); !ok {
			return false, err
		}
	}
	if gs.Precision != nil {
		if ok, err := gs.validPrecision(gs.Precision); !ok {
			return false, err
		}
	}
	if gs.TreeLevels != nil {
		if ok, err := gs.validTreeLevels(gs.TreeLevels); !ok {
			return false, err
		}
	}
	if gs.Strategy != nil {
		if ok, err := gs.validStrategy(gs.Strategy); !ok {
			return false, err
		}
	}
	if gs.DistanceErrorPct != nil {
		if ok, err := gs.validDistanceErrorPct(gs.DistanceErrorPct); !ok {
			return false, err
		}
	}
	if gs.Orientation != nil {
		if ok, err := gs.validOrientation(gs.Orientation); !ok {
			return false, err
		}
	}
	if gs.PointsOnly != nil {
		if ok, err := gs.validPointsOnly(gs.PointsOnly); !ok {
			return false, err
		}
	}

	if gs.Type == typeGeometryCollection {
		if gs.Coordinates != nil || gs.Geometries == nil {
			return false, nil
		}
	} else {
		if gs.Geometries != nil || gs.Coordinates == nil {
			return false, nil
		}
	}

	switch gs.Type {
	case typeGeometryCollection:
		str, err := piazza.StructInterfaceToString(gs.Coordinates)
		if err != nil {
			return false, err
		}
		var temp geo_GeometryCollection
		if err := json.Unmarshal([]byte(str), &temp); err != nil {
			return false, err
		}
		return temp.valid(gs)
	case typePoint:
		str, err := piazza.StructInterfaceToString(gs.Coordinates)
		if err != nil {
			return false, err
		}
		var temp geo_Sub_Point
		if err := json.Unmarshal([]byte(str), &temp); err != nil {
			return false, err
		}
		return temp.valid(gs)
	case typeLineString:
		str, err := piazza.StructInterfaceToString(gs.Coordinates)
		if err != nil {
			return false, err
		}
		var temp geo_LineString
		if err := json.Unmarshal([]byte(str), &temp); err != nil {
			return false, err
		}
		return temp.valid(gs)
	case typePolygon:
		str, err := piazza.StructInterfaceToString(gs.Coordinates)
		if err != nil {
			return false, err
		}
		var temp geo_Polygon
		if err := json.Unmarshal([]byte(str), &temp); err != nil {
			return false, err
		}
		return temp.valid(gs)
	case typeMultiPoint:
		str, err := piazza.StructInterfaceToString(gs.Coordinates)
		if err != nil {
			return false, err
		}
		var temp geo_MultiPoint
		if err := json.Unmarshal([]byte(str), &temp); err != nil {
			return false, err
		}
		return temp.valid(gs)
	case typeMultiLineString:
		str, err := piazza.StructInterfaceToString(gs.Coordinates)
		if err != nil {
			return false, err
		}
		var temp geo_MultiLineString
		if err := json.Unmarshal([]byte(str), &temp); err != nil {
			return false, err
		}
		return temp.valid(gs)
	case typeMultiPolygon:
		str, err := piazza.StructInterfaceToString(gs.Coordinates)
		if err != nil {
			return false, err
		}
		var temp geo_MultiPolygon
		if err := json.Unmarshal([]byte(str), &temp); err != nil {
			return false, err
		}
		return temp.valid(gs)
	case typeEnvelope:
		str, err := piazza.StructInterfaceToString(gs.Coordinates)
		if err != nil {
			return false, err
		}
		var temp geo_Envelope
		if err := json.Unmarshal([]byte(str), &temp); err != nil {
			return false, err
		}
		return temp.valid(gs)
	case typeCircle:
		str, err := piazza.StructInterfaceToString(gs.Coordinates)
		if err != nil {
			return false, err
		}
		var temp geo_Circle
		if err := json.Unmarshal([]byte(str), &temp); err != nil {
			return false, err
		}
		return temp.valid(gs)
	default:
		return false, nil
	}
	return true, nil
}

func (gc *geo_GeometryCollection) valid(gs *Geo_Shape) (bool, error) {
	for _, v := range *gc {
		if ok, err := v.valid(); !ok {
			return false, err
		}
	}
	return true, nil
}
func (p *geo_Sub_Point) valid(gs *Geo_Shape) (bool, error) {
	if len(*p) != 2 {
		return false, nil
	}
	for _, v := range *p {
		if /*key*/ _, ok := v.(float64); !ok {
			return false, nil
		}
	}
	point := Geo_Point{(*p)[0].(float64), (*p)[1].(float64)}
	return point.valid(), nil
}
func (ls *geo_LineString) valid(gs *Geo_Shape) (bool, error) {
	for _, v := range *ls {
		if ok, _ := v.valid(gs); !ok {
			return false, nil
		}
	}
	return true, nil
}
func (ply *geo_Polygon) valid(gs *Geo_Shape) (bool, error) {
	if len(*ply) < 1 {
		return false, nil
	}
	for _, v := range *ply {
		if len(v) != 5 {
			return false, nil
		}
		for _, p := range v {
			if ok, err := p.valid(gs); !ok {
				return false, err
			}
		}
	}
	return true, nil
}
func (mp *geo_MultiPoint) valid(gs *Geo_Shape) (bool, error) {
	for _, p := range *mp {
		if ok, err := p.valid(gs); !ok {
			return false, err
		}
	}
	return true, nil
}
func (mls *geo_MultiLineString) valid(gs *Geo_Shape) (bool, error) {
	for _, ls := range *mls {
		if ok, err := ls.valid(gs); !ok {
			return false, err
		}
	}
	return true, nil
}
func (mply *geo_MultiPolygon) valid(gs *Geo_Shape) (bool, error) {
	for _, ply := range *mply {
		if ok, err := ply.valid(gs); !ok {
			return false, err
		}
	}
	return true, nil
}
func (e *geo_Envelope) valid(gs *Geo_Shape) (bool, error) {
	if len(*e) != 2 {
		return false, nil
	}
	for _, p := range *e {
		if ok, err := p.valid(gs); !ok {
			return false, err
		}
	}
	return true, nil
}
func (c *geo_Circle) valid(gs *Geo_Shape) (bool, error) {
	p := geo_Sub_Point(*c)
	if ok, err := p.valid(gs); !ok {
		return false, err
	}
	return gs.validRadius(gs.Radius)
}

func (gs *Geo_Shape) validDistance(distance interface{}) (bool, error) {
	v, ok := distance.(string)
	if !ok {
		return false, nil
	}
	re, err := regexp.Compile(distanceRegex)
	if err != nil {
		return false, err
	}
	return re.MatchString(v), nil
}
func (gs *Geo_Shape) validRadius(radius interface{}) (bool, error) {
	return gs.validDistance(radius)
}
func (gs *Geo_Shape) validTree(tree interface{}) (bool, error) {
	v, ok := tree.(string)
	if !ok {
		return false, nil
	}
	return v == "geohash" || v == "quadtree", nil
}
func (gs *Geo_Shape) validPrecision(precision interface{}) (bool, error) {
	v, ok := precision.(string)
	if !ok {
		return ok, nil
	}
	re, err := regexp.Compile(precisionRegex)
	if err != nil {
		return false, err
	}
	return re.MatchString(v), nil
}
func (gs *Geo_Shape) validTreeLevels(treeLevels interface{}) (bool, error) {
	return gs.validDistance(treeLevels)
}
func (gs *Geo_Shape) validStrategy(strategy interface{}) (bool, error) {
	v, ok := strategy.(string)
	if !ok {
		return false, nil
	}
	return v == "recursive" || v == "term", nil
}
func (gs *Geo_Shape) validDistanceErrorPct(errorPct interface{}) (bool, error) {
	v, ok := errorPct.(float64)
	if !ok {
		return false, nil
	}
	return v >= 0 && v <= 100, nil
}
func (gs *Geo_Shape) validOrientation(orientation interface{}) (bool, error) {
	v, ok := orientation.(string)
	if !ok {
		return false, nil
	}
	re, err := regexp.Compile(orientationRegex)
	if err != nil {
		return false, err
	}
	return re.MatchString(v), nil
}
func (gs *Geo_Shape) validPointsOnly(pointsOnly interface{}) (bool, error) {
	_, ok := pointsOnly.(bool)
	return ok, nil
}