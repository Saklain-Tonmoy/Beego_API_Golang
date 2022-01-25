package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	//"log"
	"net/http"

	"github.com/beego/beego/v2/client/httplib"
	beego "github.com/beego/beego/v2/server/web"
)

type CatApiController struct {
	beego.Controller
}

type Image struct {
	Url string `json:"url"`
}

func (c *CatApiController) Index() {
	type Breed struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	}

	type Category struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
	}

	type HomeRouteData struct {
		Breed    []Breed
		Category []Category
		Image []Image
	}

	c.TplName = "index.tpl"

	/// codes for fetching all the breeds from CAT API
	/// using the "net/http" package of Golang
	breed_url := "https://api.thecatapi.com/v1/breeds?attach_breed=0"

	breed_req, _ := http.NewRequest("GET", breed_url, nil)

	breed_req.Header.Add("x-api-key", "6c73dbb1-628c-4102-b72a-cb021e2368c5")

	res, _ := http.DefaultClient.Do(breed_req)

	breed_data, _ := ioutil.ReadAll(res.Body)

	breed := []Breed{}

	json.Unmarshal(breed_data, &breed)

	/// codes for fetching all the categories from CAT API
	/// using the "httplib" package of Beego
	category_url := "https://api.thecatapi.com/v1/categories"
	catReq := httplib.Get(category_url)

	catReq.Header("Accept", "application/json")
	catReq.Header("x-api-key", "6c73dbb1-628c-4102-b72a-cb021e2368c5")

	category_data, _ := catReq.String()

	category := []Category{}

	json.Unmarshal([]byte(category_data), &category)

	fmt.Println(category)


	/// codes for fetching initial images

	imageReq := httplib.Get("https://api.thecatapi.com/v1/images/search?limit=9")

	imageReq.Header("x-api-key", "6c73dbb1-628c-4102-b72a-cb021e2368c5")

	imagesData, _ := imageReq.String()

	images := []Image{}

	json.Unmarshal([]byte(imagesData), &images)

	fmt.Println(images)


	/// forming data in structued way for serving in JSON Format
	data := HomeRouteData{}

	data.Breed = breed
	data.Category = category
	data.Image = images

	c.Data["json"] = &data

	c.ServeJSON()

}

func (c *CatApiController) GetImages() {

	order := c.GetString("order")
	mime_types := c.GetString("mime_types")
	category := c.GetString("category_ids")
	breed := c.GetString("breed_id")
	limit := c.GetString("limit")

	images_url := "https://api.thecatapi.com/v1/images/search?order=" + order + "&limit=" + limit + "&category_ids=" + category + "&breed_id=" + breed + "&mime_types=" + mime_types
	fmt.Println(images_url)

	req, _ := http.NewRequest("GET", images_url, nil)

	res, _ := http.DefaultClient.Do(req)

	images_data, _ := ioutil.ReadAll(res.Body)

	images := []Image{}

	json.Unmarshal(images_data, &images)

	c.Data["json"] = &images
	
	c.ServeJSON()
	
}
