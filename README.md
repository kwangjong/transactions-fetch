# transactions-fetch
Fetch Coding Assessment Software Engineering Internship - Backend

## Build and Run
To build and run the program, run the following commands:
```
go build
./transactions <points> <filename> --multi-thread=false
```
Or you can run the program without building, using the following command;
```
go run transactions <points> <filename> --multi-thread=false
```
\*`multi-thread=true` enables multi-threaded csv reader

## Test
To run tests, run the following command:
```
go test -v
```

## BenchMark
To run BenchMark, run the following command:
```
go test -bench=.
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

## Sample BenchMark
```
goos: darwin
goarch: amd64
pkg: transactions
cpu: Intel(R) Core(TM) i5-5250U CPU @ 1.60GHz
BenchmarkSingleThread-4   	       1	2484050592 ns/op
BenchmarkMultiThread-4    	       1	3936580385 ns/op
PASS
ok  	transactions	6.671s
```