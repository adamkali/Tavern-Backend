# Tavern-Backend
![High Resolution Logo - Transparent Background](https://user-images.githubusercontent.com/43151285/191801789-057a6cd8-9299-4cc9-97de-f140b10ac959.png)

Please see [Tavern's Readme For a proper introduction](https://github.com/adamkali/Tavern). This is just a read me for the design patterns and how to get started.

## Getting Started

If you are a part of Tavern Team, and you are not able to access the Tavern env git please email me, at the proper email, otherwise there are instructions on the link above on how to mock your own env files. Once you have everything set up, use a 
```
go clean && go build -o TavernProfile && <./TavernProfile.exe | ./TavernProfile DEPENDENT ON YOUR OS>
```

From there check to ensure that the `/api/login` and `/api/signup` and `/api/verify` endpoints all work as expected. There are examples in the Tavern git.

From there you should be able to start a branch and then submit a pull requests!

## Design Patterns

TavernProfile uses the following design pattern
```golang
// all models most comute with IData
type IData interface {
  SetID(string)
  GetID() string
  New() interface{}
}
// a model should have an ID and a Name
type Model struct {
  ID   string 
  Name string
  ...  struct { ... }
}
// a Controller for in TavernProfile should follow:
type ModelController struct{ H BaseHandler[models.Model] }

func NewUserController(DB *gorm.DB) *ModelController {
	return &ModelController{
		H: *NewHandler(DB, models.Model{}, "Model"),
	}
}
// all endpoints in the main.go should be instantiated like so
func main() {
  ...
  modelController := controllers.NewModelController(db)
  ...
  http.Handle(modelController.H.AuthPath + "/yourpath",
    cors.Handler(http.HandlerFunc(
      modelController.H.Sanitize(modelController.YourFunctionName))))
  ...
}
  // Functions that have more than one http method should use a handler with a sitch case and be called by some public function while the functions themselves should be a private function
  // in controllers/ModelController
  ...
  func ModelCustomHandler(w http.ResponseWriter, r http.Request){
    	switch r.Method {
	    case http.MethodGet:
		    getSomething(w, r)
	    case http.MethodPost:
		    h.postSomething(w, r)
	    default:
		    h.Response.UDRWrite(w, http.StatusMethodNotAllowed, "Method not allowed", false)
	    }
  }
```

As of right now there is [an issue](https://github.com/adamkali/Tavern-Backend/issues/4) that outlines adding a more practical repository outline. Check it! It will be coming in the future...

