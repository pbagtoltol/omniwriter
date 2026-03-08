# Fixed-Length to EDI Example

Demonstrates transforming fixed-length transaction records to X12 820 Payment Order/Remittance Advice EDI format using JavaScript custom functions and composite elements.

## Features Demonstrated

- **Fixed-length parsing** with `fixedlength2` format
- **EDI X12 820 output** (Payment Order/Remittance Advice)
- **JavaScript custom functions** for complex transformations
- **Composite elements** using component delimiter (:)
- **Multiple segments per transaction** (ST, N1, AMT, SE)
- **Date formatting** and decimal conversion

## Input Format

Same as fixedlength_to_csv example - fixed-length transaction file:

```
HDR20250308BATCH001
TXN0000112345SMITH     JOHN      00125099USD20250308143000
TXN0000223456JOHNSON   MARY      00087550USD20250308143100
TXN0000334567DOE       JANE      00250000CAD20250308143200
TXN0000445678WILLIAMS  ROBERT    00032575USD20250308143300
TRL0000400483224
```

## Output Format

X12 820 EDI with segments:
- **ST**: Transaction Set Header
- **N1**: Name with composite element (Last:First format)
- **AMT**: Amount with currency, amount, and date
- **SE**: Transaction Set Trailer

## JavaScript Transformations in EDI

### 1. Composite Name Element
```javascript
last.trim() + ':' + first.trim()
```
Creates composite element "SMITH:JOHN" in N1 segment.

### 2. Amount Conversion
```javascript
(cents / 100.0).toFixed(2)
```
Converts cents to decimal amount.

### 3. Date Formatting for EDI
```javascript
ts.substring(0,4) + ts.substring(4,6) + ts.substring(6,8)
```
Converts YYYYMMDDHHMMSS to YYYYMMDD for EDI date field.

## Segment Structure

Each transaction generates 4 segments:

```
ST*820*{txn_id}~                          // Transaction start
N1*{last}:{first}*{account}~              // Name (composite) + account
AMT*{currency}*{amount}*{date}~           // Amount details
SE*3*{txn_id}~                            // Transaction end (3 segments)
```

## Running the Example

```bash
go run main.go
```

Output: `output.edi`

## Expected Output

```
ST*820*00001~N1*SMITH:JOHN*12345~AMT*USD*1250.99*20250308~SE*3*00001~
ST*820*00002~N1*JOHNSON:MARY*23456~AMT*USD*875.50*20250308~SE*3*00002~
ST*820*00003~N1*DOE:JANE*34567~AMT*CAD*2500.00*20250308~SE*3*00003~
ST*820*00004~N1*WILLIAMS:ROBERT*45678~AMT*USD*325.75*20250308~SE*3*00004~
```

### Segment Breakdown (Transaction 1)

```
ST*820*00001~
  └─ Transaction Set 820, ID 00001

N1*SMITH:JOHN*12345~
  └─ Name: SMITH:JOHN (composite), Account: 12345

AMT*USD*1250.99*20250308~
  └─ Currency: USD, Amount: $1250.99, Date: 2025-03-08

SE*3*00001~
  └─ End of transaction, 3 segments, ID 00001
```

## Key Differences from CSV Example

1. **Composite Elements**: Uses colon (:) separator within N1 segment for name
2. **Multiple Segments**: Each transaction generates 4 EDI segments
3. **Date Format**: Compact YYYYMMDD instead of ISO format
4. **No Conversion**: Amount stays in original currency (no USD conversion)
5. **Segment Counting**: SE segment counts included segments (3)

## EDI Delimiters

Configured in `output_declaration`:
- Segment delimiter: `~`
- Element delimiter: `*`
- Component delimiter: `:` (used in N1 composite element)
