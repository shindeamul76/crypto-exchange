package main

import (
	"fmt"
	"sort"
	"time"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	Sizefilled float64
	Price      float64
}

type Order struct {
	Size      float64
	Bid       bool
	Limit     *Limit
	Timestamp int64
}

type Orders []*Order

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }


func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Size:      size,
		Bid:       bid,
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("[size: %.2f]", o.Size)
}

type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit

type ByBestAsks struct{ Limits }
type ByBestBids struct{ Limits }

func (a ByBestAsks) Len() int           { return len(a.Limits) }
func (a ByBestAsks) Swap(i, j int)      { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a ByBestAsks) Less(i, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }

func (b ByBestBids) Len() int           { return len(b.Limits) }
func (b ByBestBids) Swap(i, j int)      { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
func (b ByBestBids) Less(i, j int) bool { return b.Limits[i].Price > b.Limits[j].Price }



func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (l *Limit) DeleteOrder(o *Order) {
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			l.Orders[i] = l.Orders[len(l.Orders)-1]
			l.Orders = l.Orders[:len(l.Orders)-1]
		}
	}

	o.Limit = nil
	l.TotalVolume -= o.Size

	// TODO: resort the whole resting orders
	sort.Sort(l.Orders)
}

type Orderbook struct {
	Asks []*Limit
	Bids []*Limit

	AsksLimit map[float64]*Limit
	BidsLimit map[float64]*Limit
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		Asks: []*Limit{},
		Bids: []*Limit{},

		AsksLimit: make(map[float64]*Limit),
		BidsLimit: make(map[float64]*Limit),
	}
}

func (ob *Orderbook) PlaceOrder (price float64, o *Order) []Match {
	// try to match the orders with matching logic here


	// add the rest of order to the book

	if o.Size > 0.0 {
		ob.add(price, o)
	}

	return []Match{}
}

func (ob *Orderbook) add (price float64, o *Order) {

	var limit *Limit
	

	if o.Bid {
		limit = ob.BidsLimit[price]
	}else{
		limit = ob.AsksLimit[price]
	}

	if limit == nil {
		limit = NewLimit(price)
		limit.AddOrder(o)
		
		if o.Bid {
			ob.Bids = append(ob.Bids, limit)
			ob.BidsLimit[price] = limit
		}else {
			ob.Asks = append(ob.Asks, limit)
			ob.AsksLimit[price] = limit
		}
	}

}
