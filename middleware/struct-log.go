package middleware

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/link5401/fiber-chatbot/database"
	"github.com/link5401/fiber-chatbot/models"
	services "github.com/link5401/fiber-chatbot/services"
)

type LogData struct {
	Server         string `json:"server"`
	Time           string `json:"time"`
	UserName       string `json:"user_name"`
	IPAddress      string `json:"ip"`
	Method         string `json:"method"`
	Status         string `json:"status"`
	NameStatus     string `json:"names_tatus"`
	EndPoint       string `json:"end_point"`
	DataRequest    string `json:"data_request"`
	DataBeforeEdit string `json:"data_before_edit"`
	InfoResponse   string `json:"info_response"`
	InfoRequest    string `json:"info_request"`
}

func WriteLogMain(c *fiber.Ctx) error {
	// Set data log
	err := c.Next()
	services.IPAddress = c.IP()
	services.Method = c.Method()

	services.InfoRequest = string(fmt.Sprint(c.Request()))
	services.InfoResponse = string(fmt.Sprint(c.Response()))
	fmt.Println(c.Response().StatusCode())
	services.Status = fmt.Sprint(c.Response().StatusCode())

	fmt.Println(services.Status)
	services.EndPoint = c.Path()
	resp, _ := c.Request().MultipartForm()
	if resp != nil {
		data, _ := c.Request().MultipartForm()
		services.DataRequest = fmt.Sprint(data)
	} else {
		if string(c.Request().URI().QueryString()) != "" { // type GET
			services.DataRequest = string(c.Request().URI().QueryString())
		} else {
			if string(c.Body()) != "" {
				services.DataRequest = string(covertJsonBodyToString(c))
			} else {
				services.DataRequest = ""
			}
		}
	}

	if services.Status != "200" {
		InsertLog(c)
	}
	return err
}

func WriteLogMidleWare(c *fiber.Ctx) {
	// Set data log

	services.IPAddress = c.IP()
	services.Method = c.Method()
	services.Status = fmt.Sprint(fiber.StatusOK)
	services.EndPoint = c.Path()

	services.InfoRequest = string(fmt.Sprint(c.Request()))
	services.InfoResponse = string(fmt.Sprint(c.Response()))

	resp, _ := c.Request().MultipartForm()
	if resp != nil {
		data, _ := c.Request().MultipartForm()
		services.DataRequest = fmt.Sprint(data)
	} else {
		if string(c.Request().URI().QueryString()) != "" { // type GET
			services.DataRequest = string(c.Request().URI().QueryString())
		} else {
			if string(c.Body()) != "" {
				services.DataRequest = string(covertJsonBodyToString(c))
			} else {
				services.DataRequest = ""
			}
		}
	}

	if services.Status != "200" {
		InsertLog(c)
	}

}

// Write Log Error Controller
// @ Error Code
// @ c *fiber.Ctx
func WriteLogErrorController(errCode string, c *fiber.Ctx) error {
	// convert DataBeforeEditStruct(struct -> string)
	if services.DataBeforeEditStruct != nil {
		services.DataBeforeEdit = covertJsonToString(services.DataBeforeEditStruct)
	}
	services.Status = errCode
	InsertLog(c)

	//set data response is null
	services.DataBeforeEditStruct = nil
	services.DataBeforeEdit = ""
	return c.JSON(fiber.Map{
		"message": "error",
	})
}

// Write Log Success Controller
// @ DataBeforeEdit
// @ c *fiber.Ctx
func WriteLogSuccessController(c *fiber.Ctx) error {
	// convert DataBeforeEditStruct(struct -> string)
	if services.DataBeforeEditStruct != nil {
		services.DataBeforeEdit = covertJsonToString(services.DataBeforeEditStruct)
	}
	InsertLog(c)

	//set data response is null
	services.DataBeforeEditStruct = nil
	services.DataBeforeEdit = ""

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func covertJsonToString(Json interface{}) string {
	out, _ := json.Marshal(Json)
	personMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(out), &personMap)
	if err != nil {
		panic(err)
	}
	var strOut = "["
	for key, value := range personMap {
		if fmt.Sprint(value) == "<nil>" {
			value = ""
		}
		strOut += key + ":[" + fmt.Sprint(value) + "] "
	}
	var str = fmt.Sprint(strings.Trim(string(strOut), " ") + "]")
	return str
}

func covertJsonBodyToString(c *fiber.Ctx) string {
	personMap := make(map[string]interface{})
	err := json.Unmarshal(c.Body(), &personMap)
	if err != nil {
		panic(err)
	}
	var strOut = "["
	for key, value := range personMap {
		if fmt.Sprint(value) == "<nil>" {
			value = ""
		}
		strOut += key + ":[" + fmt.Sprint(value) + "] "
	}
	var str = fmt.Sprint(strings.Trim(string(strOut), " ") + "]")
	return str
}

////////////write log to db/////////////
func InsertLog(c *fiber.Ctx) error {
	var err error
	db := database.DB
	tbl_log := new(models.LogCustom)

	dataRequestCovert := strings.Replace(services.DataRequest, "\u0026{map", "", 1)
	NameStatus := ""
	if services.Status == "200" {
		NameStatus = "success"
	} else {
		NameStatus = "error"
	}

	tbl_log.Time = time.Now().Format("2006-01-02 15:04:05")
	tbl_log.UserName = services.UserName
	tbl_log.IPAddress = services.IPAddress
	tbl_log.Method = services.Method
	tbl_log.Status = services.Status
	tbl_log.NameStatus = NameStatus
	tbl_log.EndPoint = services.EndPoint
	tbl_log.DataRequest = strings.Replace(dataRequestCovert, " map[]}", "", 1)
	tbl_log.DataBeforeEdit = services.DataBeforeEdit

	request := strings.Split(services.InfoRequest, "------WebKitFormBoundary")
	infoRequest := strings.Split(request[0], "----------------------------")

	tbl_log.InfoRequest = string(infoRequest[0])
	tbl_log.InfoResponse = string(services.InfoResponse)

	errors := models.ValidateLogCustom(*tbl_log)

	if errors != nil {
		return c.Status(400).JSON(errors)
	}

	err = db.Create(&tbl_log).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "system_error",
		})
	}

	// Return response
	return c.JSON(fiber.Map{
		"message": "success",
	})
}
