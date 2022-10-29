package check

// Catch returns err if calling work panics with Error{err}, otherwise it
// returns nil.
//
//	return check.Catch(func() {
//		check.Must1(fmt.Println("Hello, World!")
//		check.Must1(fmt.Println("¡Hola, Mundo!")
//		check.Must1(fmt.Println("你好，世界!")
//		check.Must1(fmt.Println("Привет, мир!")
//	})
func Catch(work func(), transforms ...func(e error) error) (e error) {
	defer Handle(&e, transforms...)
	work()
	return
}

// Catch1 returns _, err if calling work panics with Error{err}, otherwise it
// returns t, nil.
//
//	func getTotalWeight(weight, qty string) (float64, error) {
//		return Catch1(func() float64 {
//			return Must1(strconv.ParseFloat(weight, 64)) *
//				float64(Must1(strconv.Atoi(qty)))
//		})
//	}
func Catch1[T any](work func() T, transforms ...func(e error) error) (t T, e error) {
	defer Handle(&e, transforms...)
	t = work()
	return
}

// Catch2 returns _, _, err if calling work panics with Error{err}, otherwise it
// returns t1, t2, nil. See Catch1 for a related example.
func Catch2[T1, T2 any](
	work func() (T1, T2),
	transforms ...func(e error) error,
) (t1 T1, t2 T2, e error) {
	defer Handle(&e, transforms...)
	t1, t2 = work()
	return
}

// Catch4 returns _, _, err if calling work panics with Error{err}, otherwise it
// returns t1, t2, t3, nil. See Catch1 for a related example.
func Catch3[T1, T2, T3 any](
	work func() (T1, T2, T3),
	transforms ...func(e error) error,
) (t1 T1, t2 T2, t3 T3, e error) {
	defer Handle(&e, transforms...)
	t1, t2, t3 = work()
	return
}

// Catch4 returns _, _, _, err if calling work panics with Error{err}, otherwise
// it returns t1, t2, t3, t4 nil. See Catch1 for a related example.
func Catch4[T1, T2, T3, T4 any](
	work func() (T1, T2, T3, T4),
	transforms ...func(e error) error,
) (t1 T1, t2 T2, t3 T3, t4 T4, e error) {
	defer Handle(&e, transforms...)
	t1, t2, t3, t4 = work()
	return
}
