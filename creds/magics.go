package creds

func unmagic(cr []byte) {
	for i := 0; i < len(cr); i++ {
		cr[i]--
	}
}

func magic(cr []byte) {
	for i := 0; i < len(cr); i++ {
		cr[i]++
	}
}
