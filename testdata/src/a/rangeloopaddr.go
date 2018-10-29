package a

func RangeLoopAddrIntTests() {
	var s []int

	// No loop variables.
	for range s {
	}
	for _ = range s {
	}
	for _, _ = range s {
	}

	for i, j := range s {
		// Legal access to values.
		println(i, j)

		// Illegal access to addresses.
		println(&i, j)  // want "taking address of range variable 'i'"
		println(i, &j)  // want "taking address of range variable 'j'"
		println(&i, &j) // want "taking address of range variable 'i'" "taking address of range variable 'j'"

		// Legal access to shadowing variables.
		i := i
		j := j
		println(&i)
		println(&j)
	}

	// Wrap in dummy range to catch erroneous early bailout.
	for range s {
		for i := range s {
			println(&i) // want "taking address of range variable 'i'"
		}
	}

	{
		var i, j int

		for i = range s {
			println(i)
			println(&i) // want "taking address of range variable 'i'"
		}

		for _, j = range s {
			println(j)
			println(&j) // want "taking address of range variable 'j'"
		}

		for i, j = range s {
			println(i, j)
			println(&i) // want "taking address of range variable 'i'"
			println(&j) // want "taking address of range variable 'j'"
		}
	}

	for i := range s {
		func(i int) {
			// Legal access of shadowing variable.
			println(&i)
		}(i)
	}
	{
		var i int
		ip := &i
		for i = range s {
			println(ip)
		}
	}

	func() *int {
		for i, j := range s {
			// Legal access in return statement.
			for range s {
				return &i
			}
			for range s {
				return &j
			}
			for i, _ := range s {
				return &i
			}
			for _, j := range s {
				return &j
			}

			// Illegal access in return statement of nested function.
			func() *int {
				return &i // want "taking address of range variable 'i'"
			}()

			func() *int {
				return &j // want "taking address of range variable 'j'"
			}()
			func() *int {
				for range s {
					return &j // want "taking address of range variable 'j'"
				}
				return nil
			}()

			func() *int {
				for i := range s {
					// Legal access when in return statement.
					return &i
				}
				for i := range s {
					println(&i) // want "taking address of range variable 'i'"
				}
				for range s {
					return &j // want "taking address of range variable 'j'"
				}
				return nil
			}()

			// Legal access in return statement.
			return &i
		}
		return nil
	}()

	func() func() *int {
		for i := range s {
			return func() *int {
				// Legal access in return statement in defining loop.
				for range s {
					println(&i)
					return &i
				}
				println(&i)
				return &i
			}
		}
		return nil
	}()

	func() func() *int {
		for i := range s {
			return func() *int {
				for i := range s {
					// Illegal access in shadowing variable.
					println(&i) // want "taking address of range variable 'i'"
				}
				for i := range s {
					// Legal access in of return statement of shadowing variable.
					return &i
				}
				println(&i)
				for range s {
					// Legal access in nested dummy loop.
					println(&i)
				}
				for j := range s {
					// Legal access of parent variable but illegal access to the one in nested loop.
					println(&i, &j) // want "taking address of range variable 'j'"
				}
				return &i
			}
		}
		return nil
	}()
}

func RangeLoopAddrStructTests() {
	var vs []struct {
		m int
	}

	for _, v := range vs {
		// Illegal access to addresses of members.
		println(&v.m) // want "taking address of range variable 'v'"
	}
}
