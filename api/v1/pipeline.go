package v1

import (
	"encoding/json"
	guuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/klovercloud-ci/api/common"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/api"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type pipelineApi struct {
	pipelineService service.Pipeline
}
var (
	upgrader = websocket.Upgrader{}
)
func (p pipelineApi) GetLogs(context echo.Context)error {
	processId:=context.Param("processId")
	option := getQueryOption(context)
	logs,total:=p.pipelineService.GetLogsByProcessId(processId,option)
	metadata := common.GetPaginationMetadata(option.Pagination.Page, option.Pagination.Limit, total, int64(len(logs)))
	uri:=strings.Split(context.Request().RequestURI,"?")[0]
	if option.Pagination.Page > 0 {
		metadata.Links = append(metadata.Links, map[string]string{"prev": uri + "?order=" + context.QueryParam("order") + "&page=" + strconv.FormatInt(option.Pagination.Page-1, 10) + "&limit=" + strconv.FormatInt(option.Pagination.Limit, 10)})
	}
	metadata.Links = append(metadata.Links, map[string]string{"self": uri+ "?order=" + context.QueryParam("order") + "&page=" + strconv.FormatInt(option.Pagination.Page, 10) + "&limit=" + strconv.FormatInt(option.Pagination.Limit, 10)})

	if (option.Pagination.Page+1)*option.Pagination.Limit < metadata.TotalCount {
		metadata.Links = append(metadata.Links, map[string]string{"next": uri + "?order=" + context.QueryParam("order") + "&page=" + strconv.FormatInt(option.Pagination.Page+1, 10) + "&limit=" + strconv.FormatInt(option.Pagination.Limit, 10)})
	}
	return common.GenerateSuccessResponse(context,logs,&metadata,"")
}

func getQueryOption(context echo.Context) v1.LogEventQueryOption {
	option := v1.LogEventQueryOption{}
	page := context.QueryParam("page")
	limit := context.QueryParam("limit")
	if page == "" {
		option.Pagination.Page = enums.DEFAULT_PAGE
		option.Pagination.Limit = enums.DEFAULT_PAGE_LIMIT
	} else {
		option.Pagination.Page, _ = strconv.ParseInt(page ,10, 64)
		option.Pagination.Limit, _ = strconv.ParseInt(limit ,10, 64)
	}
	return option
}

func (p pipelineApi) GetEvents(context echo.Context) error {
	processId:=context.QueryParam("processId")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(context.Response(), context.Request(), nil)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer ws.Close()

	status :=make(chan map[string]interface{})
	for {
		go p.pipelineService.ReadEventByProcessId(status,processId)
		jsonStr, err := json.Marshal(<-status)
		if err!=nil{
			log.Println(err.Error())
		}
		err = ws.WriteMessage(websocket.TextMessage, []byte(jsonStr))
		if err != nil {
			context.Logger().Error(err)
		}
		_, _, err = ws.ReadMessage()
		if err != nil {
			context.Logger().Error(err)
		}

	}
}

func (p pipelineApi) Apply(context echo.Context) error {
	var data v1.Pipeline
	body, err := ioutil.ReadAll(context.Request().Body)
	if  err != nil{
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context,nil,err.Error())
	}
	if err := json.Unmarshal(body, &data); err != nil {
		if err := yaml.Unmarshal(body, &data); err != nil {
			return common.GenerateErrorResponse(context,nil,err.Error())
		}
	}
	url:=context.QueryParam("url")
	revision:=context.QueryParam("revision")
	purgingOption:=context.QueryParam("purging")
	if purgingOption==string(enums.PIPELINE_PURGING_ENABLE){
		data.Option.Purging=enums.PIPELINE_PURGING_ENABLE
	}else{
		data.Option.Purging=enums.PIPELINE_PURGING_DISABLE
	}
	data.ApiVersion =string(enums.API_V1)
	if data.ProcessId!=""{
		data.ProcessId = guuid.New().String()
	}
	err = p.pipelineService.Apply(url,revision, data)

	if err != nil{
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context,err.Error(),"Failed to trigger pipeline!")
	}
	return common.GenerateSuccessResponse(context,data.ProcessId,nil,"Pipeline successfully triggered!")
}

func NewPipelineApi(pipelineService service.Pipeline) api.Pipeline {
	return &pipelineApi{
		pipelineService: pipelineService,
	}
}
