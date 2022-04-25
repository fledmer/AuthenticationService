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
	RoleID      string  `bson:"roleID"`
	Address     Address `bson:"address"`
	Verified    bool    `bson:"verified"`
}

type UserRoles struct {
	ID   string `bson:"ID"`
	Name string `bson:"Name"`
}

type Deed struct {
	ID          string   `bson:"ID"`
	CreatorID   string   `bson:"creatorID"`
	Description string   `bson:"description"`
	Categories  string   `bson:"categories"`
	Topic       string   `bson:"topic"`
	Address     Address  `bson:"address"`
	ImagesId    []string `bson:"imagesId"`
}

type Images struct {
	ID   string `bson:"ID"`
	Path string `bson:"path"`
	Url  string `bson:"url"`
}

type CookieSession struct {
	ID    string `bson:"ID"`
	Token string `bson:"token"`
}
