Geometry
========

####A 2D geometry library for Go, based on postgres datatypes.

Allows the following postgres datatypes to be sent and received by enhanced postgres database/sql drivers.  Currently the only driver which supports this is my fork of lib/pq, found at github.com/gregb/pq.

The following postgres datatypes are supported:

|----|----|
| *Postgres Type* | *Go Type* |
|----|----|
|point|geometry.Point|
|point|geometry.Vector|
|line||
|lseg|geometry.Segment|
|box|geometry.Box|
|circle|geometry.Circle|
|path*||
|polygon*||


From: http://www.postgresql.org/docs/9.3/static/datatype-geometric.html
