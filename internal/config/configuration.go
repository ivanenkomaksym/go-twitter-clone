package config

type ApiServer struct {
	ApplicationUrl string
}

type TweetsStorage struct {
	ConnectionString string
	DatabaseName     string
}

type FeedsStorage struct {
	ConnectionString    string
	DatabaseName        string
	FeedsCollectionName string
}

type Configuration struct {
	ApiServer     ApiServer
	TweetsStorage TweetsStorage
	FeedsStorage  FeedsStorage
}
