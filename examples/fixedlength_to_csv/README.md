# Fixed-Length to CSV Example

Demonstrates transforming fixed-length transaction records to CSV format using complex JavaScript custom functions.

## Features Demonstrated

- **Fixed-length parsing** with `fixedlength2` format
- **JavaScript custom functions** for data transformation:
  - String concatenation and formatting
  - Conditional logic (currency conversion)
  - Date/time formatting
  - Number formatting with `toFixed()`
- **Multiple envelopes** (header, transaction, trailer)
- **CSV output** with custom columns

## Input Format

Fixed-length transaction file with three record types:
- **HDR** (Header): Batch information
- **TXN** (Transaction): Customer payment records
- **TRL** (Trailer): Summary totals

```
HDR20250308BATCH001
TXN0000112345SMITH     JOHN      00125099USD20250308143000
TXN0000223456JOHNSON   MARY      00087550USD20250308143100
...
TRL0000400483224
```

### Field Positions

Transaction records (TXN):
- Position 1-3: Record type ("TXN")
- Position 4-8: Transaction ID
- Position 9-13: Account number
- Position 14-23: Last name (padded)
- Position 24-33: First name (padded)
- Position 34-41: Amount in cents
- Position 42-44: Currency code
- Position 45-58: Timestamp (YYYYMMDDHHmmss)

## Output Format

CSV with 7 columns:
- transaction_id: Formatted as "TXN-0000000001"
- customer_name: "Last, First" format
- account: Account number
- amount_usd: Converted to USD (CAD * 0.75)
- currency_code: Original currency
- transaction_date: ISO format "YYYY-MM-DD HH:MM:SS"
- description: Generated description

## JavaScript Transformations

### 1. Transaction ID Formatting
```javascript
'TXN-' + txn.padStart(10, '0')
```
Pads ID with leading zeros.

### 2. Name Concatenation
```javascript
last.trim() + ', ' + first.trim()
```
Combines and trims names.

### 3. Currency Conversion
```javascript
curr === 'USD' ? (cents / 100.0).toFixed(2) : (cents / 100.0 * 0.75).toFixed(2)
```
Converts CAD to USD using 0.75 rate.

### 4. Timestamp Parsing
```javascript
ts.substring(0,4) + '-' + ts.substring(4,6) + '-' + ts.substring(6,8) + ' ' +
ts.substring(8,10) + ':' + ts.substring(10,12) + ':' + ts.substring(12,14)
```
Converts YYYYMMDDHH mmss to ISO format.

### 5. Description Generation
```javascript
'Payment from account ' + acct + ' in ' + curr
```
Creates descriptive text.

## Running the Example

```bash
go run main.go
```

Output: `output.csv`

## Expected Output

```csv
TXN-0000000001,"SMITH, JOHN",12345,1250.99,USD,2025-03-08 14:30:00,Payment from account 12345 in USD
TXN-0000000002,"JOHNSON, MARY",23456,875.50,USD,2025-03-08 14:31:00,Payment from account 23456 in USD
TXN-0000000003,"DOE, JANE",34567,1875.00,CAD,2025-03-08 14:32:00,Payment from account 34567 in CAD
TXN-0000000004,"WILLIAMS, ROBERT",45678,325.75,USD,2025-03-08 14:33:00,Payment from account 45678 in USD
```

Note: Row 3 shows CAD conversion (2500.00 CAD → 1875.00 USD)
