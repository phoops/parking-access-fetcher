## Description

A job to take off-street parking access data from Kafka and insert it into Scorpio broker using mobility toolkit data model Vehicle. Plate numbers are removed.

NOTE: The coordinates of Vehicles on Scorpio Broker are sent inverted. This is due to the fact that mobility toolkit uses inverted coordinates (probably is a bug).

## Run

`go run cmd/nurse/main.go`

## License

for the ODALA project.

© 2023 Phoops

License EUPL 1.2

![CEF Logo](https://ec.europa.eu/inea/sites/default/files/ceflogos/en_horizontal_cef_logo_2.png)

The contents of this publication are the sole responsibility of the authors and do not necessarily reflect the opinion of the European Union.
This project has received funding from the European Union’s “The Connecting Europe Facility (CEF) in Telecom” programme under Grant Agreement number: INEA/CEF/ICT/A2019/2063604
