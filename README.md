# 💉 Nurse

## Description
A job to take parking data from kafka and insert it into scorpio broker in mobility toolkit data model Vehicle. Plate number has been removed

NOTE: The coordinates of Vehicles on Scorpio Broker are saved inverted. This is due to the fact that mobility toolkit uses inverted coordinates (probably is a bug).


## Run
`go run cmd/nurse/main.go`