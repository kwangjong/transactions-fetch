# transactions-fetch
Fetch Coding Assessment Software Engineering Internship - Backend

## Build and Run
To build and run the program, run the following commands:
```
go build
./transaction <points> <filename>
```
Or you can run the program without building, using the following command;
```
go run transactions <points> <filename>
```

## Test
To run test, run the following command:
```
go test -v
```

## Sample Input
```
[transactions.csv]
"payer","points","timestamp"
"DANNON",1000,"2020-11-02T14:00:00Z"
"UNILEVER",200,"2020-10-31T11:00:00Z"
"DANNON",-200,"2020-10-31T15:00:00Z"
"MILLER COORS",10000,"2020-11-01T14:00:00Z"
"DANNON",300,"2020-10-31T10:00:00Z"
```

## Sample Output
```
$ go run transactions 5000 transactions.csv 
{
        "DANNON": 1000,
        "MILLER COORS": 5300,
        "UNILEVER": 0
}
```

