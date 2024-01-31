## 0. prerequisite

- uuid : file_id를 생성하기 위해 사용
- gin : http server를 구성하기 위해 사용
- gorm : db를 구성하기 위해 사용
- postgres : db를 구성하기 위해 사용

### environment

#### docker & postgres

```shell
docker run --name easywalk_postgres -e POSTGRES_PASSWORD=easywalk -e POSTGRES_DB=easywalk -p 5432:5432 -d postgres
```

### uuid, gin, gorm, gorm postgres, postgres, easywalk

```shell
go get -u github.com/google/uuid
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/lib/pq
go get -u github.com/easywalk/go-restful
```

# 1. file-api 구현

### model 생성 (pk/model/model.go)

```go
type File struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
    Name      string    `json:"name" gorm:"type:varchar(255);not null;unique;"`
    Size      uint64    `json:"size" gorm:"type:bigint;"`
    Type      string    `json:"type" gorm:"type:varchar(255);"`
    CreatedAt time.Time `json:"createdAt" gorm:"type:timestamp;"`
    UpdatedAt time.Time `json:"updatedAt" gorm:"type:timestamp;"`
}

func (f *File) SetCreatedAt(t time.Time) {
    f.CreatedAt = t
}

func (f *File) SetUpdatedAt(t time.Time) {
    f.UpdatedAt = t
}

func (f *File) GetID() uuid.UUID {
    return f.ID
}

func (f *File) SetID(id uuid.UUID) {
    f.ID = id
}

```

### easywalk를 이용한 http server 구성 (main.go)

```go
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
```

# 2. 테스트

```http request
### 생성
### 생성
< {%
client.global.clearAll()
request.variables.clearAll()
request.variables.set("base_name", "test");
%}
POST localhost:8080/files
Content-Type: application/json

{
    "name": "{{base_name}}-{{$uuid}}",
    "size": 100,
    "type": "image/jpeg"
}
> {%
client.test("Status code is 200", function () {
client.assert(response.status === 200, "Response status is not 200");
client.assert(response.body.ID !== "", "Response body ID is empty");

client.global.set("id", response.body.id);
client.global.set("name", response.body.name);
client.global.set("size", response.body.size);
client.global.set("type", response.body.type);
});
%}

### 생성 검증
GET localhost:8080/files/{{id}}
Content-Type: application/json

> {%
client.test("Status code is 200", function () {
client.assert(response.status === 200, "Response status is not 200");
client.assert(response.body.name === client.global.get("name"), "Response body Name is not same :" + response.body.name + " : " + client.global.get("name"));
client.assert(response.body.size == client.global.get("size"), "Response body Size is not same :" + response.body.size + " : " + client.global.get("size"));
client.assert(response.body.type === client.global.get("type"), "Response body Type is not same :" + response.body.type + " : " + client.global.get("type"));
});
%}


### 재생성 (실패 - 이름 중복)
POST localhost:8080/files
Content-Type: application/json

{
    "name": "{{name}}"
}
> {%
client.test("Status code is 500", function () {
client.assert(response.status === 500, "Response status is not 200");
});
%}

### 이름 수정
PATCH localhost:8080/files/{{id}}
Content-Type: application/json

{
    "name": "{{base_name}}-{{$uuid}}"
}

> {%
client.test("Status code is 200", function () {
client.assert(response.status === 200, "Response status is not 200");
client.assert(response.body.id === client.global.get("id"), "Response body ID is not same :" + response.body.ID + " : " + client.global.get("id"));
client.assert(response.body.name !== client.global.get("name"), "Response body Name is not same :" + response.body.name + " : " + client.global.get("name"));
client.assert(response.body.size == client.global.get("size"), "Response body Size is not same :" + response.body.size + " : " + client.global.get("size"));
client.assert(response.body.type === client.global.get("type"), "Response body Type is not same :" + response.body.type + " : " + client.global.get("type"));

client.global.set("id", response.body.id);
client.global.set("name", response.body.name);
client.global.set("size", response.body.size);
client.global.set("type", response.body.type);
});
%}

### 이름 수정 검증
GET localhost:8080/files/{{id}}
Content-Type: application/json

> {%
client.test("Status code is 200", function () {
client.assert(response.status === 200, "Response status is not 200");
client.assert(response.body.id === client.global.get("id"), "Response body ID is not same :" + response.body.ID + " : " + client.global.get("id"));
client.assert(response.body.name === client.global.get("name"), "Response body Name is not same :" + response.body.name + " : " + client.global.get("name"));
client.assert(response.body.size == client.global.get("size"), "Response body Size is not same :" + response.body.size + " : " + client.global.get("size"));
client.assert(response.body.type === client.global.get("type"), "Response body Type is not same :" + response.body.type + " : " + client.global.get("type"));
});
%}

### 삭제
DELETE localhost:8080/files/{{id}}
Content-Type: application/json

> {%
client.test("Status code is 204", function () {
client.assert(response.status === 204, "Response status is not 204" + response.status);
});
%}

### 없는 것 또 삭제 (idempotent)
DELETE localhost:8080/files/{{id}}
Content-Type: application/json

> {%
client.test("Status code is 204", function () {
client.assert(response.status === 204, "Response status is not 204" + response.status);
});
%}

### 삭제 검증
GET localhost:8080/files/{{id}}
Content-Type: application/json

> {%
client.test("Status code is 204", function () {
client.assert(response.status === 204, "Response status is not 404" + response.status);
});
%}
```