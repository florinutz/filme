package list

type Paging struct {
	filterLow    int
	filterHigh   int
	limitLow     int
	limitHigh    int
	itemsPerPage int
	pagesToCrawl []int
}

// getNextPages returns the range of pages that should be scheduled for parsing.
// Should only be called after a page was crawled.
func (p Paging) getNextPages(maxItems int) (pages []int) {
	if p.filterLow > p.limitHigh {
		return
	}

	wantedPages := (maxItems-1)/p.itemsPerPage + 1

	if p.filterHigh > 0 && p.filterLow > 0 {
		wantedPages = p.filterHigh - p.filterLow + 1
	}

	if p.limitHigh < wantedPages+1 {
		wantedPages = p.limitHigh - 1
	}

	bottom := p.limitLow
	top := bottom + wantedPages - 1

	if p.filterLow > bottom {
		bottom = p.filterLow
		top = bottom + wantedPages - 1
	}

	if p.filterHigh > 0 && p.filterHigh < top {
		top = p.filterHigh
	}

	if bottom > top {
		return
	}

	for i := bottom; i <= top; i++ {
		pages = append(pages, i)
	}

	return
}

func (p Paging) pageIsValid(pageNo, wantedItems int) bool {
	if wantedItems == 0 && pageNo == 1 {
		return true
	}

	neededPages := p.getNextPages(wantedItems)
	if len(neededPages) == 0 {
		return false
	}
	return pageNo >= neededPages[0] && pageNo <= neededPages[len(neededPages)-1]
}
