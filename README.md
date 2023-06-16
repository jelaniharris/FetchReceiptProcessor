# Receipt Processor

This is a webservice that allows you to create a receipt and calculate a point value of that receipt using a set of rules.

## How to run

### Start the webservice
```
go run .
```

### Running the tests
```
go test -v ./...
```

## Endpoints

### View All Receipts
* Path: `/receipts`
* Method: `GET`

View all of the receipts in the system

### View Receipt

* Path: `/receipts/{id}`
* Method: `GET`

View the contents of a receipt

### Process Receipt

* Path: `/receipts/process`
* Method: `POST`

Takes a receipt via JSON and then returns the id of the created receipt

### Calculate Points
* Path: `/receipts/{id}/points`
* Method: `GET`

This endpoint takes a receipt id and returns the number of points that receipt awarded

The rules are as follows:

#### Rules

* One point for every alphanumeric character in the retailer name.
  * "Target" is woth 6 points
  * "Barnes & Noble" is woth 11 points (spaces and symbols don't count)
* 50 points if the total is a round dollar amount with no cents.
  * "10.00" is worth 50, "4.50" is worth 0
* 25 points if the total is a multiple of `0.25`.
  * "4.00" is worth 25, "3.16" is worth 0
* 5 points for every two items on the receipt.
  * So a list with 4 items has two pairs with 10 points
  * A list with 9 items has four pairs and is worth 20 points
* If the trimmed length of the item description is a multiple of 3, multiply the price by `0.2` and round up to the nearest integer. The result is the number of points earned.
  * "  Target  " has a length of 6, which is a multiple of 3, so the points would be PRICE * 0.2
* 6 points if the day in the purchase date is odd.
  * "2023-06-15" would be worth 6 points because 15 is an odd numbered day
* 10 points if the time of purchase is after 2:00pm and before 4:00pm.
  * "15:40" is 3:40PM on a 24-hour clock so it would count for 10 points
