package tableview

type Pagination struct {
	nextConfirmPending bool
	prevConfirmPending bool
	message            string
}

func (p *Pagination) HasPendingConfirm() bool {
	return p.nextConfirmPending || p.prevConfirmPending
}

func (p *Pagination) Message() string {
	return p.message
}

func (p *Pagination) RequestNextPage() {
	p.nextConfirmPending = true
	p.prevConfirmPending = false
	p.message = "Press G again to load next page"
}

func (p *Pagination) RequestPreviousPage() {
	p.prevConfirmPending = true
	p.nextConfirmPending = false
	p.message = "Press g again to load previous page"
}

func (p *Pagination) ConfirmNextPage() bool {
	if p.nextConfirmPending {
		p.Clear()
		return true
	}
	return false
}

func (p *Pagination) ConfirmPrevPage() bool {
	if p.prevConfirmPending {
		p.Clear()
		return true
	}
	return false
}

func (p *Pagination) RequestPrevPage() {
	p.prevConfirmPending = true
	p.nextConfirmPending = false
	p.message = "Press g again to load previous page"
}

func (p *Pagination) Clear() {
	p.nextConfirmPending = false
	p.prevConfirmPending = false
	p.message = ""
}
