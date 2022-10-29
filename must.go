package check

// Must calls panic(Error{err}) if err is not nil.
func Must(err error) {
	if err != nil {
		panic(Error{err})
	}
}

// Must1 returns t if err is nil, otherwise it calls panic(Error{err}).
//
//	price := check.Must1(strconv.ParseFloat(unitPrice, 64)) *
//		check.Must1(strconv.ParseFloat(qty, 64))
func Must1[T any](t T, err error) T {
	if err != nil {
		panic(Error{err})
	}
	return t
}

// Must2 returns t1, t2 if err is nil, otherwise it calls panic(Error{err}).
//
//	// MulDiv's third return value is an error if x = y = 0.
//	prod, quo := check.Must2(MulDiv(x, y))
func Must2[T1, T2 any](t1 T1, t2 T2, err error) (T1, T2) {
	Must(err)
	return t1, t2
}

// Must3 returns t1, t2, t3 if err is nil, otherwise it calls panic(Error{err}).
//
//	// MulDivRem's fourth return value is an error if x = y = 0.
//	prod, quo, rem := check.Must3(MulDivRem(x, y))
func Must3[T1, T2, T3 any](t1 T1, t2 T2, t3 T3, err error) (T1, T2, T3) {
	Must(err)
	return t1, t2, t3
}

// Must4 returns t1, t2, t3, t4 if err is nil, otherwise it calls panic(Error{err}).
//
//	// AnalyzeTrades's fifth return value is an error if prices is empty.
//	open, high, low, close := check.Must4(AnalyzeTrades(prices))
func Must4[T1, T2, T3, T4 any](t1 T1, t2 T2, t3 T3, t4 T4, err error) (T1, T2, T3, T4) {
	Must(err)
	return t1, t2, t3, t4
}
