package crossref

import "github.com/caltechlibrary/crossrefapi"

type ApiRequestItem interface {
	Execute(*crossrefapi.CrossRefClient)
}

type ApiResponse[T crossrefapi.WorksQueryResponse | crossrefapi.Works] struct {
	Data  *T
	Error error
}

type QueueItem[T crossrefapi.WorksQueryResponse | crossrefapi.Works] struct {
	Query    func(client *crossrefapi.CrossRefClient) ApiResponse[T]
	Receiver chan ApiResponse[T]
}

func (i QueueItem[T]) Execute(client *crossrefapi.CrossRefClient) {
	i.Receiver <- i.Query(client)
}

type Queue[T crossrefapi.WorksQueryResponse | crossrefapi.Works] chan QueueItem[T]

func NewQueue[T crossrefapi.WorksQueryResponse | crossrefapi.Works](cap int) Queue[T] {
	return make(Queue[T], cap)
}

func (q Queue[T]) Push(item QueueItem[T]) {
	q <- item
}

func (q Queue[T]) Pop() (*QueueItem[T], bool) {
	select {
	case item := <-q:
		return &item, true
	default:
		return nil, false
	}
}

func (q Queue[T]) Len() int {
	return len(q)
}
