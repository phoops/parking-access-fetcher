# 💉 Nurse

## Description
A job to take off-street parking access data from Kafka and insert it into Scorpio broker using mobility toolkit data model Vehicle. Plate numbers are removed.

NOTE: The coordinates of Vehicles on Scorpio Broker are sent inverted. This is due to the fact that mobility toolkit uses inverted coordinates (probably is a bug).


## Run
`go run cmd/nurse/main.go`
