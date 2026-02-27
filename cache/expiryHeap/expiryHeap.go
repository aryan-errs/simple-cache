package expiryheap

import "aryan-errs/simple-cache/cache/types"

type ExpiryHeap []*types.ExpiryItem

func (h ExpiryHeap) Len() int {
	return len(h)
}

func (h ExpiryHeap) Less(i int, j int) bool {
	return h[i].ExpiresAt.Before(h[j].ExpiresAt)
}

func (h ExpiryHeap) Swap(i int, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h ExpiryHeap) Peak() types.ExpiryItem {
	return *h[0]
}

func (h *ExpiryHeap) Push(x any) {
	*h = append(*h, x.(*types.ExpiryItem))
}

func (h *ExpiryHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = &types.ExpiryItem{}
	*h = old[0 : n-1]
	return x
}
