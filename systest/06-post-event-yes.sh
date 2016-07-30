#!/bin/bash

source setup.sh

#url="http://pz-workflow.$PZDOMAIN"

eventTypeId=$1
if [ "$eventTypeId" == "" ]
then
    echo "error: \$eventTypeId missing"
    exit 1
fi

json='{
    "eventTypeId": "'"$eventTypeId"'",
    "data": {
        "beta": 17,
        "alpha": "quick brown fox"
    }
}'

#echo POST /event
#echo "$json"

ret=$($curl -XPOST -d "$json" "$url"/event)

#echo RETURN:
#echo "$ret"

echo `extract eventId "$ret"`
