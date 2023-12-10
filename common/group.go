package common

// GroupList 将长度为max的数据切分开来，每部分batch个，最后不够batch有多少算多少
// var tt = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
// ch := GroupList(len(tt), 3)
//
//	for k := range ch {
//	    fmt.Println(k, tt[k[0]:k[1]])
//	}
func GroupList(max, batch int) <-chan [2]int {
	ch := make(chan [2]int, max/batch+1)
	defer close(ch)
	for i := 0; i <= max-1; i += batch {
		end := i + batch
		if max < end {
			end = max
		}
		ch <- [2]int{i, end}
	}
	return ch
}
