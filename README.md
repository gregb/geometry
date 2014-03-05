Geometry
========

####A 2D geometry library for Go, based on postgres datatypes.

Allows the following postgres datatypes to be sent and received by enhanced postgres database/sql drivers.  Currently the only driver which supports this is my fork of lib/pq, found at http://github.com/gregb/pq.

The following postgres datatypes are supported:

| Postgres Type | Go Type |
| ---------- | ----------|
| point| geometry.Point |
| point| geometry.Vector |
| line | *(1) Not supported |
| lseg | geometry.Segment |
| box | geometry.Box |
| circle | geometry.Circle |
| path  | *(2) Not yet supported |
| polygon | *(2) Not yet supported |

From: http://www.postgresql.org/docs/9.3/static/datatype-geometric.html

*(1) Postgres docs indicate "line" support is incomplete.  And since I could not distinguish any use cases distinct from "segment", I did not implement it.

*(2) I've had no personal need for these, and have therefore not bothered to write anything.  I'd be happy to consider pull requests.
