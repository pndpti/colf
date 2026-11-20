# colf file format 

## layout 

```
["COLF" magic, 4 bytes][version, 1 byte]
[column 0 data][column 1 data] ...
[footer bytes]
[footer length, uint32 little-endian]
["COLF" magic, 4 bytes]
```