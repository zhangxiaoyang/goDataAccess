package main

func test() []int {
	println("C")
	return []int{1, 2, 3}
}
func main() {
	for _, v := range test() {
		println(v)
	}
}
