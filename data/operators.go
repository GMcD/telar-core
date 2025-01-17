package data

type SearchOperator struct {
	Search string `json:"$search" bson:"$search"`
}

type TextOperator struct {
	Text SearchOperator `json:"$text" bson:"$text"`
}

type UpdateOperator struct {
	Set interface{} `json:"$set" bson:"$set"`
}

type AddToSetOperator struct {
	AddToSet interface{} `json:"$addToSet" bson:"$addToSet"`
}

type IncrementOperator struct {
	Inc interface{} `json:"$inc" bson:"$inc"`
}
