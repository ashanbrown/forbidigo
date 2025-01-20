package gokeyword

func Do1() error {
	return nil
}

func Do2() {
	if err := Do1(); err != nil {
		panic(err) // want "forbidden by pattern"
	}
}

func Do3() {
	go Do2() // want "forbidden by pattern"
}
