## layout

```
["COLF" magic, 4 bytes][version, 1 byte]
[column 0 data][column 1 data] ...
[footer bytes]
[footer length, uint32]
["COLF" magic, 4 bytes]
```

## column data

- int64: each value as 8 bytes
- string: each value as uvarint length + bytes

## footer

```
uvarint numRows
uvarint numCols
per column:
  uvarint nameLen + name bytes
  1 byte type (0 = int64, 1 = string)
  uvarint offset   (from start of file)
  uvarint size     (bytes)
```
