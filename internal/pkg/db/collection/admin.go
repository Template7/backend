package collection

type Admin struct {
	//Id       *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string              `json:"username" bson:"username" binding:"required"`
	Password string              `json:"password" bson:"password" binding:"required"`
}
