package auth

import (
	"github.com/spf13/viper"
	"testing"
	"time"
)

func TestRefreshToken(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	viper.Set("Mongo.Db", "testDb")

	//testUser := structs.User{
	//	Mobile: "+886987654321",
	//}
	//uId, err := user.CreateUser(testUser)
	//if err != nil {
	//	t.Error(err)
	//}
	//token, err := GenUserToken(uId.Hex())
	//if err != nil {
	//	t.Error(err)
	//}
	////db.New().SaveToken()
	//
	//type args struct {
	//	oriToken structs.Token
	//}
	//tests := []struct {
	//	name               string
	//	args               args
	//	wantRefreshedToken structs.Token
	//	wantErr            *t7Error.Error
	//}{
	//	{
	//		name: "normal",
	//		args: args{
	//			oriToken: token,
	//		},
	//		wantErr: nil,
	//	},
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		_, gotErr := RefreshToken(tt.args.oriToken)
	//		if !reflect.DeepEqual(gotErr, tt.wantErr) {
	//			t.Errorf("RefreshToken() gotErr = %v, want %v", gotErr, tt.wantErr)
	//		}
	//	})
	//}
	//
	//c, _ := mongo.Connect(nil, options.Client().ApplyURI(config.New().Mongo.ConnectionString))
	//if err != nil {
	//	t.Error(err)
	//}
	//db := c.Database(config.New().Mongo.Db)
	//db.Drop(context.Background())
}

func TestTemp(t *testing.T) {
	t.Log("start...1")

	time.Sleep(1 * time.Second)

	f := func() {
		t.Log("start go routine...2")
		defer func() {
			t.Log("defer before sleep...3")
			time.Sleep(1 * time.Second)
			t.Log("defer after sleep...4")
		}()
	}

	f()

	t.Log("...5")

	time.Sleep(2 * time.Second)

	t.Log("test end...6")
}
