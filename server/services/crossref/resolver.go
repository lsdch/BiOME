package crossref

type CrossRefResolver struct {
	client *CrossRefScheduler
}

func NewCrossRefResolver(client *CrossRefScheduler) *CrossRefResolver {
	return &CrossRefResolver{
		client: client,
	}
}
