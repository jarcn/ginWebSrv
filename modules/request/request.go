package request

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func GetJson(c *gin.Context) (map[string]interface{}, error) {
	jsonStr, _ := ioutil.ReadAll(c.Request.Body)
	var data map[string]interface{}
	err := json.Unmarshal(jsonStr, &data)
	return data, err
}
