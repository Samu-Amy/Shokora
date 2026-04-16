package main

// ----- FUNCTIONS -----

func randomFrom[T any](arr []T) T {
	return arr[customRand.Intn(len(arr))]
}

func randomBool() bool {
	return customRand.Float32() < 0.5
}
