#!/bin/bash
INDEX_NAME=alerts004
ALIAS_NAME=$1
ES_IP=$2
TESTING=$3

AlertMapping='
	"Alert": {
		"dynamic": "strict",
		"properties": {
			"alertId": {
				"type": "string",
				"index": "not_analyzed"
			},
			"triggerId": {
				"type": "string",
				"index": "not_analyzed"
			},
			"jobId": {
				"type": "string",
				"index": "not_analyzed"
			},
			"eventId": {
				"type": "string",
				"index": "not_analyzed"
			},
			"createdBy": {
				"type": "string",
				"index": "not_analyzed"
			},
			"createdOn": {
				"type": "date",
				"format": "yyyy-MM-dd'\''T'\''HH:mm:ssZZ||yyyy-MM-dd'\''T'\''HH:mm:ss.SZZ||yyyy-MM-dd'\''T'\''HH:mm:ss.SSZZ||yyyy-MM-dd'\''T'\''HH:mm:ss.SSSZZ||yyyy-MM-dd'\''T'\''HH:mm:ss.SSSSZZ||yyyy-MM-dd'\''T'\''HH:mm:ss.SSSSSZZ||yyyy-MM-dd'\''T'\''HH:mm:ss.SSSSSSZZ||yyyy-MM-dd'\''T'\''HH:mm:ss.SSSSSSSZZ"
			}
		}
	}'

IndexSettings="
{
	"\""mappings"\"": {
		$AlertMapping
	}
}"


bash db/CreateIndex.sh $INDEX_NAME $ALIAS_NAME $ES_IP "$IndexSettings" "$AlertMapping" $TESTING