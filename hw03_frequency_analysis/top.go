package hw03frequencyanalysis

import (
	"bufio"
	"container/heap"
	"strings"
)

func countWords(str string) map[string]uint64 {
	wordCount := map[string]uint64{}

	scanner := bufio.NewScanner(strings.NewReader(str))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		wordCount[word]++
	}

	return wordCount
}

type wordStat struct {
	word  string
	count uint64
}

type wordStatHeap []wordStat

func (w wordStatHeap) Len() int { return len(w) }
func (w wordStatHeap) Less(i, j int) bool {
	return w[i].count > w[j].count || (w[i].count == w[j].count && w[i].word < w[j].word)
}
func (w wordStatHeap) Swap(i, j int) { w[i], w[j] = w[j], w[i] }

func (w *wordStatHeap) Push(x any) {
	*w = append(*w, x.(wordStat))
}

func (w *wordStatHeap) Pop() any {
	n := len(*w)
	x := (*w)[n-1]
	*w = (*w)[0 : n-1]
	return x
}

func Top10(str string) []string {
	const numOfTop = 10

	wordCount := countWords(str)

	h := &wordStatHeap{}
	for word, count := range wordCount {
		heap.Push(h, wordStat{word, count})
	}

	top := make([]string, 0, min(h.Len(), numOfTop))

	for i := 0; i < numOfTop && h.Len() > 0; i++ {
		top = append(top, heap.Pop(h).(wordStat).word)
	}

	return top
}
