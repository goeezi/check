// Package check implements an exception-handling system for Go using panic and
// recover under the hood, with generics enabling a fairly clean API.
//
// Because generics don't offer variadic type parameter packs, package check
// provides a family of Catch and Must functions for up to four explicitly
// defined parameter types.
//
// # Example
//
// The following example shows a piece of code written in Go's conventional
// error handling approach.
//
//	type Data struct {
//	    conn             *sql.DB
//	    selectPricesStmt *sql.Stmt
//	}
//
//	func New(driver, dsn string) (*Data, error) {
//	    conn, err := sql.Open(driver, dsn)
//
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    return &Data{conn: conn}, nil
//	}
//
//	func (d *Data) GetPrices(symbol string) (open, hi, lo, close float64, err error) {
//	    if d.selectPricesStmt == nil {
//	        q, err := d.conn.Prepare("SELECT o, h, l, c FROM prices WHERE symbol = ?")
//	        if err != nil {
//	            return 0, 0, 0, 0, err
//	        }
//	        d.selectPricesStmt = q
//	    }
//
//	    tx, err := d.conn.Begin()
//	    if err != nil {
//	        return 0, 0, 0, 0, err
//	    }
//
//	    return getPrices(tx.Stmt(d.selectPricesStmt), symbol)
//	}
//
//	func getPrices(stmt *sql.Stmt, symbol string) (open, hi, lo, close float64, err error) {
//	    q, err := stmt.Query()
//	    if err != nil {
//	        return 0, 0, 0, 0, err
//	    }
//	    colTypes, err := q.ColumnTypes()
//	    if err != nil {
//	        return 0, 0, 0, 0, err
//	    }
//	    log.Printf("column types: %#v", colTypes)
//	    if q.Next() {
//	        err := q.Scan(&open, &hi, &lo, &close)
//	        if err != nil {
//	            return 0, 0, 0, 0, err
//	        }
//	        if q.Next() {
//	            return 0, 0, 0, 0, fmt.Errorf("too many results for %q", symbol)
//	        }
//	    } else {
//	        return 0, 0, 0, 0, fmt.Errorf("no results for %q", symbol)
//	    }
//	    return
//	}
//
// Below is the above example reworked using package check.
//
//	type Data struct {
//	    conn             *sql.DB
//	    selectPricesStmt *sql.Stmt
//	}
//
//	func New(driver, dsn string) (_ *Data, err error) {
//	    defer check.Handle(&err)
//	    return &Data{conn: check.Must1(sql.Open(driver, dsn))}, nil
//	}
//
//	func (d *Data) GetPrices(symbol string) (open, hi, lo, close float64, err error) {
//	    return check.Catch4(func() (_, _, _, _ float64) {
//	        defer check.Handle(&err)
//	        if d.selectPricesStmt == nil {
//	            d.selectPricesStmt = check.Must1(
//	                d.conn.Prepare("SELECT o, h, l, c FROM prices WHERE symbol = ?"))
//	        }
//
//	        tx := check.Must1(d.conn.Begin())
//
//	        return getPrices(tx.Stmt(d.selectPricesStmt), symbol)
//	    })
//	}
//
//	func getPrices(stmt *sql.Stmt, symbol string) (open, hi, lo, close float64) {
//	    q := check.Must1(stmt.Query())
//	    log.Printf("column types: %#v", check.Must1(q.ColumnTypes()))
//	    if q.Next() {
//	        check.Must(q.Scan(&open, &hi, &lo, &close))
//	        if q.Next() {
//	            check.Fail(fmt.Errorf("too many results for %q", symbol))
//	        }
//	    } else {
//	        check.Fail(fmt.Errorf("no results for %q", symbol))
//	    }
//	    return
//	}
//
// The most obvious difference between the two examples is the length of the
// code. The original weighs in at 36 significant lines of code, while the
// second is just 24, a reduction of exactly 1/3.
//
// A less obvious, but more significant distinction is a reduction from eight
// internal variables (ignoring input and return parameters) down to just two.
// This represents a sharp drop in the amount of internal state and a matching
// reduction in the amount of mental bookkeeping required to comprehend the flow
// of logic. As an example, the colTypes variable is used four lines after it is
// defined in the original code, and the experienced reader is predisposed to
// keep a mental note of it after that, just in case it crops up later in the
// code. They might even wonder, "Why is it here? Is it just for logging or does
// the function have some other use for it?" While the reader might be barely
// (or not even) aware of these thoughts, they will nonetheless clutter the mind
// as the logic increases in scope and complexity. In the rewritten example, the
// colTypes variable doesn't exist at all. The expression is used directly, which
// doesn't trigger any of the above questions, and the reader, instinctively
// knowing that it won't be referred to again, can simply discard that sliver of
// information. In fact, they will likely skim past the log.Printf call without
// even being consciously aware of it.
//
// Several other points are worth noting:
//
//  1. There are two available methods for trapping errors before they escape
//     a package: check.Handle and check.Catch/CatchN. The first uses a
//     conventional deferred function call to recover a panicked error and
//     convert it back to a returned error. The second uses the same technique
//     internally, but calls a function to do the work. This second form is less
//     idiomatic, but has the advantage of being usable in any block scope, not
//     just function-level scope.
//
//  2. Not every function must trap errors. Note that the unpublished getPrices
//     function uses check.Must/MustN, but doesn't use check.Handle or
//     check.Catch/CatchN. This is perfectly acceptable usage within a package,
//     since the published methods will trap errors before they escape.
//
//  3. MustN and CatchN, only go up to 4 parameters. To deal with functions that
//     return more than four return values plus an error, assign their output to
//     local variables the conventional way then call check.Must(err).  In
//     practice, one should generally not create functions with more than four
//     return values plus an error. They can invariably be redesigned to return
//     a struct.
//
//  5. All instances of returning default values have disappeared in the new
//     code. This is another important way in which package check reduces
//     cognitive load, both on the author and the reader.
//
// # Performance considerations
//
// From the profile below, it is clear that error handling using package check
// is much slower than conventional error handling.
//
//	❯ go test -run=^$ -bench=. -benchmem -cpuprofile cpu.success.out
//	goos: darwin
//	goarch: arm64
//	pkg: github.com/anzx/acceleration-tools/envelope/cmd/envelope/internal/check
//	BenchmarkFailureConventional-8          1000000000           0.3117 ns/op          0 B/op          0 allocs/op
//	BenchmarkFailureCatch-8                  6531588           180.7 ns/op        16 B/op          1 allocs/op
//	BenchmarkFailureHandle-8                 8418494           140.8 ns/op        16 B/op          1 allocs/op
//	BenchmarkFailureHandleTransform-8        8462744           143.1 ns/op        16 B/op          1 allocs/op
//	BenchmarkSuccessConventional-8          1000000000           0.3106 ns/op          0 B/op          0 allocs/op
//	BenchmarkSuccessCatch-8                 140923567            8.558 ns/op           0 B/op          0 allocs/op
//	BenchmarkSuccessHandle-8                240914712            5.008 ns/op           0 B/op          0 allocs/op
//	BenchmarkSuccessHandleTransform-8       200517948            5.921 ns/op           0 B/op          0 allocs/op
//	PASS
//	ok      github.com/anzx/acceleration-tools/envelope/cmd/envelope/internal/check 13.524s
//
// Conventional error handling clocks in at just over 0.3 ns regardless of
// whether the call succeeds or fails.
//
// In contrast, a successful call to check.Handle is almost 20 times slower and
// almost 30 times slower when calling check.Catch.
//
// Things are much worse during failures. Failed calls to check.Handle and
// check.Catch are 500 and 600 times slower, respectively, than conventional
// error handling.
//
// The clear message from this analysis is to avoid package check
// in performance sensitive code. That said, it is worth keeping things in
// perspective. A 5–8 ns overhead for successful calls is still very fast and
// would be perfectly acceptable in most contexts. More thought would need to be
// given to scenarios where errors are common, but even then a failed call
// still takes a small fraction of the time it takes to perform most forms of
// I/O.
package check
