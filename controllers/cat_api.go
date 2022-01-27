package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type CatApiController struct {
	beego.Controller
}

type Image struct {
	Url string `json:"url"`
}

func FetchApi(url string, ch chan string) {

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatalln(err)
	}

	response, _ := http.DefaultClient.Do(request)

	data, _ := ioutil.ReadAll(response.Body)

	ch <- string(data)

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
		Image    []Image
	}


	breedChannel := make(chan string)
	categoryChannel := make(chan string)
	imageChannel := make(chan string)

	go FetchApi("https://api.thecatapi.com/v1/breeds", breedChannel)

	go FetchApi("https://api.thecatapi.com/v1/categories", categoryChannel)

	go FetchApi("https://api.thecatapi.com/v1/images/search?limit=9", imageChannel)

	breeds := []Breed{}
	json.Unmarshal([]byte(<-breedChannel), &breeds)
	//fmt.Println(breed)

	categories := []Category{}
	json.Unmarshal([]byte(<-categoryChannel), &categories)
	//fmt.Println(category)

	images := []Image{}
	json.Unmarshal([]byte(<-imageChannel), &images)
	//fmt.Println(image)

	close(breedChannel)
	close(categoryChannel)
	close(imageChannel)

	c.TplName = "index.tpl"

	c.Data["breeds"] = breeds
	c.Data["categories"] = categories
	c.Data["images"] = images


	data := HomeRouteData{}

	data.Breed = breeds
	data.Category = categories
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
	page := c.GetString("page")

	images_url := "https://api.thecatapi.com/v1/images/search?order=" + order + "&limit=" + limit + "&category_ids=" + category + "&breed_id=" + breed + "&mime_types=" + mime_types + "&page=" + page
	fmt.Println(images_url)

	req, _ := http.NewRequest("GET", images_url, nil)

	res, _ := http.DefaultClient.Do(req)

	images_data, _ := ioutil.ReadAll(res.Body)

	images := []Image{}

	json.Unmarshal(images_data, &images)

	c.Data["json"] = &images

	c.ServeJSON()

}
