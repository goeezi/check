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
// error handling approach on the left diffed with a version using package check
// on the right.
//
//	type Data struct {                            · type Data struct {
//	  db        *sql.DB                           ·   db        *sql.DB
//	  selPrices *sql.Stmt                         ·   selPrices *sql.Stmt
//	}                                             · }
//	                                              ·
//	func New(driver, dsn string) (*Data, error) { |  func New(driver, dsn string) (_ *Data, e error) {
//	  conn, err := sql.Open(driver, dsn)          |    defer check.Handle(&e)
//	  if err != nil {                             <
//	    return nil, err                           <
//	  }                                           <
//	  return &Data{db: conn}, nil                 |    return &Data{
//	                                              >      db: check.Must1(sql.Open(driver, dsn)),
//	                                              >    }, nil
//	}                                             · }
//	                                              ·
//	func (d *Data) GetPrices(sym string) (        · func (d *Data) GetPrices(sym string) (
//	  open, hi, lo, close float64, _ error,       |   open, hi, lo, close float64, e error,
//	) {                                           · ) {
//	                                              >   defer check.Handle(&e)
//	  if d.selPrices == nil {                     ·   if d.selPrices == nil {
//	    q, err := d.db.Prepare(                   |     d.selPrices = check.Must1(d.db.Prepare(
//	      `SELECT o,h,l,c                         ·       `SELECT o,h,l,c
//	       FROM price                             ·        FROM price
//	       WHERE sym=?`)                          |        WHERE sym=?`))
//	    if err != nil {                           <
//	      return 0, 0, 0, 0, err                  <
//	    }                                         <
//	    d.selPrices = q                           <
//	  }                                           ·   }
//	  tx, err := d.db.Begin()                     |   tx := check.Must1(d.db.Begin())
//	  if err != nil {                             <
//	    return 0, 0, 0, 0, err                    <
//	  }                                           <
//	  return getPrices(tx.Stmt(d.selPrices), sym) |   o, h, l, c := getPrices(tx.Stmt(d.selPrices), sym)
//	                                              >   return o, h, l, c, nil
//	}                                             · }
//	                                              ·
//	func getPrices(stmt *sql.Stmt, sym string) (  · func getPrices(stmt *sql.Stmt, sym string) (
//	  open, hi, lo, close float64, _ error,       |   open, hi, lo, close float64,
//	) {                                           · ) {
//	  q, err := stmt.Query()                      |   q := check.Must1(stmt.Query())
//	  if err != nil {                             <
//	    return 0, 0, 0, 0, err                    <
//	  }                                           <
//	  colTypes, err := q.ColumnTypes()            |   log.Printf("cols: %#v", check.Must1(q.ColumnTypes()
//	  if err != nil {                             <
//	    return 0, 0, 0, 0, err                    <
//	  }                                           <
//	  log.Printf("cols: %#v", colTypes)           <
//	  if q.Next() {                               ·   if q.Next() {
//	    err := q.Scan(&open, &hi, &lo, &close)    |     check.Must(q.Scan(&open, &hi, &lo, &close))
//	    if err != nil {                           <
//	      return 0, 0, 0, 0, err                  <
//	    }                                         <
//	    if q.Next() {                             ·     if q.Next() {
//	      return 0, 0, 0, 0, fmt.Errorf("> 1 resu…|       check.Failf("> 1 result: %q", sym)
//	    }                                         ·     }
//	  } else {                                    ·   } else {
//	    return 0, 0, 0, 0, fmt.Errorf("no result:…|     check.Failf("no result: %q", sym)
//	  }                                           ·   }
//	  return                                      ·   return
//	}                                             · }
//
// The most obvious difference between the two examples is the length of the
// code. The original weighs in at 39 significant lines of code, while the
// second is just 29, a reduction of 25%.
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
// colTypes variable doesn't exist at all. The expression is used directly,
// which doesn't trigger any of the above questions, and the reader,
// instinctively knowing that it won't be referred to again, can simply discard
// that sliver of information. In fact, they will likely skim past the
// log.Printf call without even being consciously aware of it.
//
// Several other points are worth noting:
//
//  1. Not every function must trap errors. Note that the unpublished getPrices
//     function uses check.Must/MustN, but doesn't use check.Handle or
//     check.Catch/CatchN. This is perfectly acceptable usage within a package,
//     since the published methods will trap errors before they escape.
//
//  2. MustN and CatchN only go up to 4 parameters. To deal with functions that
//     return more than four return values plus an error, assign their output to
//     local variables the conventional way then call check.Must(err).  In
//     practice, one should generally not create functions with more than four
//     return values plus an error. They are usually better redesigned to return
//     a struct.
//
//  3. All instances of returning "don't care" zero values have disappeared in
//     the new code. This is another important way in which package check
//     reduces cognitive load, both on the author and the reader.
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
// The clear message from this analysis is to avoid using package check in
// performance sensitive code. That said, it is worth keeping things in
// perspective. A 5–8 ns overhead for successful calls is still very fast and
// would be perfectly acceptable in most contexts. More thought would need to be
// given to scenarios where errors are common, but even then a failed call still
// takes a small fraction of the time it takes to perform most forms of I/O.
package check
