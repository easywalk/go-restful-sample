package main

import (
	"github.com/easywalk/go-restful-sample/pkg/model"
	"github.com/easywalk/go-restful/easywalk/handler"
	"github.com/easywalk/go-restful/easywalk/repository"
	"github.com/easywalk/go-restful/easywalk/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	dbUserName := "postgres"
	dbPassword := "easywalk"
	dbName := "easywalk"
	dsn := "host=localhost user=" + dbUserName + " password=" + dbPassword + " dbname=" + dbName + " port=5432 sslmode=disable TimeZone=Asia/Seoul"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	r := gin.Default()
	group := r.Group("/files")
	// create File Service
	repo := repository.NewSimplyRepository[*model.File](db)
	svc := service.NewGenericService[*model.File](repo)
	hdlr := handler.NewHandler[*model.File](group, svc)
	if hdlr != nil {
		log.Println("Success to create File Handler")
	}

	r.Run() // listen and serve on
}
