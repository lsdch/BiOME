package services

// ApiRequestItem represents an item in the API request queue that can be executed.
// The generic parameter Client defines the type of client instance needed for execution.
type ApiRequestItem[Client any] interface {
	Execute(*Client)
}

// ApiResponse represents a generic API response structure that can hold both data and error information.
type ApiResponse[T any] struct {
	Data  *T
	Error error
}

// QueueItem represents a generic queue item that contains a query function and a response channel.
// It's designed to handle asynchronous API requests where T is the expected response type
// and Client is the client type used to make the request.
type QueueItem[T any, Client any] struct {
	Query    func(client *Client) ApiResponse[T]
	Receiver chan ApiResponse[T]
}

func NewQueueItem[T, Client any](
	fn func(client *Client) ApiResponse[T],
) QueueItem[T, Client] {
	return QueueItem[T, Client]{
		Query:    fn,
		Receiver: make(chan ApiResponse[T]),
	}
}

// Execute processes the queue item using the provided client.
func (i QueueItem[T, Client]) Execute(client *Client) {
	i.Receiver <- i.Query(client)
}

// Queue is a generic channel type that handles QueueItem elements.
type Queue[T any, Client any] chan QueueItem[T, Client]

func NewQueue[T any, Client any](cap int) Queue[T, Client] {
	return make(Queue[T, Client], cap)
}

func (q Queue[T, Client]) Push(item QueueItem[T, Client]) {
	q <- item
}

func (q Queue[T, Client]) Pop() (*QueueItem[T, Client], bool) {
	select {
	case item := <-q:
		return &item, true
	default:
		return nil, false
	}
}

func (q Queue[T, Client]) Len() int {
	return len(q)
}
