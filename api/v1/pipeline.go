package v1

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/klovercloud-ci/api/common"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/api"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	"github.com/labstack/echo/v4"
	guuid "github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
)

type pipelineApi struct {
	pipelineService service.Pipeline
}
var (
	upgrader = websocket.Upgrader{}
)
func (p pipelineApi) GetLog(context echo.Context)error {
	processId:=context.Param("processId")
	option := getQueryOption(context)
	logs,total:=p.pipelineService.GetLogsByProcessId(processId,option)
	metadata := common.GetPaginationMetadata(option.Pagination.Page, option.Pagination.Limit, total, int64(len(logs)))
	if option.Pagination.Page > 0 {
		metadata.Links = append(metadata.Links, map[string]string{"prev": context.Path() + "?order=" + context.QueryParam("order") + "&page=" + strconv.FormatInt(option.Pagination.Page-1, 10) + "&limit=" + strconv.FormatInt(option.Pagination.Limit, 10)})
	}
	metadata.Links = append(metadata.Links, map[string]string{"self": context.Path() + "?order=" + context.QueryParam("order") + "&page=" + strconv.FormatInt(option.Pagination.Page, 10) + "&limit=" + strconv.FormatInt(option.Pagination.Limit, 10)})

	if (option.Pagination.Page+1)*option.Pagination.Limit < metadata.TotalCount {
		metadata.Links = append(metadata.Links, map[string]string{"next": context.Path() + "?order=" + context.QueryParam("order") + "&page=" + strconv.FormatInt(option.Pagination.Page+1, 10) + "&limit=" + string(option.Pagination.Limit)})
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
	log.Println(processId)
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(context.Response(), context.Request(), map[string][]string{"Access-Control-Allow-Origin": {"*"}} )
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer ws.Close()
	//status :=make(chan string)
	for {
		//go p.pipelineService.ReadEventByProcessId(processId)
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint(p.pipelineService.ReadEventByProcessId(processId))))
		if err != nil {
			context.Logger().Error(err)
		}
		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			context.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}

func (p pipelineApi) Apply(context echo.Context) error {
	data:=v1.Pipeline{}
	err := context.Bind(&data)
	if  err != nil{
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context,nil,err.Error())
	}
	url:=context.QueryParam("url")
	revision:=context.QueryParam("revision")

	data.ApiVersion = enums.Api_version
	data.ProcessId = guuid.New().String()
	error := p.pipelineService.Apply(url,revision, data)

	if error != nil{
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
