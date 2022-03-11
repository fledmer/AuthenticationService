package systemComponents

type Address struct {
	StreetType string `bson:"streetType"`
	StreetName string `bson:"streetName"`
	House      string `bson:"house"`
	Corps      string `bson:"corps"`
	Apartment  string `bson:"apartment"`
}

type User struct {
	ID          string  `bson:"ID"`
	Mail        string  `bson:"mail"`
	Password    string  `bson:"password"`
	FirstName   string  `bson:"firstName"`
	LastName    string  `bson:"lastName"`
	ThirdName   string  `bson:"thirdName"`
	Birthday    string  `bson:"birthday"`
	PhoneNumber string  `bson:"phoneNumber"`
	Age         uint32  `bson:"age"`
	IsAdmin     bool    `bson:"isAdmin"`
	Address     Address `bson:"address"`
}

type Deed struct {
	ID          string  `bson:"ID"`
	CreatorID   string  `bson:"creatorID"`
	Description string  `bson:"description"`
	Categories  string  `bson:"categories"`
	Topic       string  `bson:"topic"`
	Address     Address `bson:"address"`
}
