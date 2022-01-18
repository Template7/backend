package user

import (
	"fmt"
	"github.com/Template7/backend/internal/pkg/config"
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	viper.AddConfigPath("../../../configs")
	c := config.New()
	db := fmt.Sprintf("temp_test")
	c.Mongo.Db = db
	c.MySql.Db = db
	c.MySql.ConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.MySql.Username, c.MySql.Password, c.MySql.Host, c.MySql.Port, c.MySql.Db)
	code := m.Run()

	//teardown(db)
	os.Exit(code)
}

func Test_user(t *testing.T) {
	//CreateUser()
}
