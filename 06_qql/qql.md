# Introduction

QQL is a non-relational database language

# DML

## SET TABLE

QQL works with non-relational tables, every command assumes a table is _set_ or _active_
as the current table being worked on.

To set a table active, use:

```
set TABLENAME
```

## SELECT

Given a currently active table, you can select data with an empty statement,
that is equivalent to `select * from active_table`, without a `where` clause.

Any subsequent expressions can be used to filter the data:

```
name ~ "artur" & bday = "1999.09.26"
```

The above expression is equivalent to `select * from active_table where name like '%artur%' and bday = '1999.09.26'`.

## INSERT

```
+ name="Artur" bday="1999.09.26"
```

The above expression is equivalent to `insert into active_table (name, bday) values ('Artur', '1999.09.26')`, since
the table is known, and the columns are ordered, this can be shortened to: `+ 'Artur' '1999.09.26'`

## DELETE

```
- name~"Artur" bday="1999.09.26"
```
The above expression is equivalent to `delete from active_table where name like '%Artur%' and bday = '1999.09.26'`

## UPDATE

```
! name~"Artur" bday="1999.09.26" -> name="Guilherme" bday="2000.04.02"
```

The above expression is equivalent to `update active_table set name = 'Guilherme', bday='2000.04.02' where name like '%Artur%' and bday = '1999.09.26'`.
Can also be shortened to:

```
! name~"Artur" bday="1999.09.26" -> "Guilherme" "2000.04.02"
```

# DDL

## DATA TYPES

Types are built from values, there are 3 categories for values:
- integer: `1`, `2`, `0x1`, `0b1101`
- float: `1.0e-10`, `2.0`, `1.0e2`
- rune: `'a'`, `'B'`

There are a two operators on these values that create new types:
- `|` means alternation, `1 | 0` means the value is either the integer `1` or the integer `0`.
You can only use alternation on the same value category, `'a' | 1` is forbidden.
- Justaposition means concatenation, `1 2` means the value is equivalent to `1` followed by `2`

Additionally, there are a few syntax sugars:
- `..` is the range operator: `1 .. 3` is the same thing as `1 | 2 | 3`
- `[i]` is the array operator: `[3]1` is the same thing as `1 1 1`
- `"abc"` is the string operator: `"abc"` is the same thing as `'a' 'b' 'c'`

You can assign names to types, but types can't be recursive in any way.

Examples:
```
type Byte     {0 .. 255}
type Float32  {1.175494351e-38 .. 3.402823466e+38}
type Status   {"Done" | "Canceled" | "Pending"}
type Letter   {'a' .. 'z' | 'A' .. 'Z' }
type NameChar {Letter | ' '}
type Name     {[45]NameChar}
type CodeName {Letter Letter Letter}

type Num {'0' .. '9'}
type CPF {Num Num Num '.' Num Num Num '.' Num Num Num '-' Num Num}

type IsoDate {Num Num Num Num '.' Num Num '.' Num Num}
type Year {-10000 .. 10000}
type Month {1 .. 12}
type Day {1 .. 31}
type StrucDate { Year Month Day }

type RGB {Byte Byte Byte}

type Hex {'0' .. '9' | 'A' .. 'F'}
type HexRGB {Hex Hex Hex Hex Hex Hex}
```

## CREATE TABLE

```
create TABLENAME -> name:Name bday:IsoDate
```

## ALTER TABLE

### Add Column
```
alter TABLENAME -> + height:Int
```
### Remove Column
```
alter TABLENAME -> - height
```
### Alter Column
```
alter TABLENAME -> ! height -> height:Int
```

## DROP TABLE
```
drop TABLENAME
```
