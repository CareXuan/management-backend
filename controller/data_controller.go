package controller

import (
	"bytes"
	"data_verify/common"
	"data_verify/conf"
	"data_verify/model"
	"data_verify/service"
	"fmt"
	dbf "github.com/SebastiaanKlippert/go-foxpro-dbf"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// GBKDecoder 自定义解码器
type GBKDecoder struct{}

// Decode 实现 dbf.Decoder 接口的 Decode 方法
func (d *GBKDecoder) Decode(b []byte) ([]byte, error) {
	// 使用 GBK 编码进行解码
	decoded, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(b), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, err
	}
	return decoded, nil // 返回 []byte
}

// UploadSbkDBF 上传文件并读取内容
func UploadSbkDBF(c *gin.Context) {
	//year, err := common.GetCurrentYear()
	//if err != nil {
	//	common.ResError(c, "获取年份信息失败")
	//	return
	//}
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		common.ResError(c, "获取文件失败")
		return
	}

	// 将文件保存到临时目录
	tempFile, err := ioutil.TempFile("", "uploaded-*.dbf")
	if err != nil {
		common.ResError(c, "创建文件失败")
		return
	}
	defer os.Remove(tempFile.Name()) // 处理完成后删除临时文件

	if err := c.SaveUploadedFile(file, tempFile.Name()); err != nil {
		common.ResError(c, "保存文件失败")
		return
	}

	// 打开 DBF 文件并使用 GBK 解码器
	dbfTable, err := dbf.OpenFile(tempFile.Name(), &GBKDecoder{})
	if err != nil {
		common.ResError(c, "读取文件失败")
		return
	}
	defer dbfTable.Close()

	var bkdItems []*model.CheckData
	sqlErr := conf.Mysql.Find(&bkdItems)
	if sqlErr != nil {
		common.ResError(c, "获取报考点失败")
		return
	}
	var bkdMapping = make(map[string]map[string][]string)
	for _, i := range bkdItems {
		bkdMapping[i.Bmddm] = make(map[string][]string)
		bkdMapping[i.Bmddm]["kmdy2"] = []string{}
		bkdMapping[i.Bmddm]["kmdy3"] = []string{}
		bkdMapping[i.Bmddm]["kmdy4"] = []string{}
	}

	// 遍历并输出记录内容
	for !dbfTable.EOF() {
		record, err := dbfTable.Record()
		if err != nil {
			common.ResError(c, "读取记录失败")
			return
		}
		dbfTable.Skip(1)

		// 跳过已删除的记录
		if record.Deleted {
			continue
		}
		ksbh := strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("KSBH"))).(string), " ", "")
		bmddm := strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("BMDDM"))).(string), " ", "")
		wgym := strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("WGYM"))).(string), " ", "")
		ywk1m := strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("YWK1M"))).(string), " ", "")
		ywk2m := strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("YWK2M"))).(string), " ", "")
		if strings.HasPrefix(wgym, "20") && wgym != "-" {
			bkdMapping[bmddm]["kmdy2"] = append(bkdMapping[bmddm]["kmdy2"], ksbh)
		}
		if strings.HasPrefix(ywk1m, "30") && wgym != "-" {
			bkdMapping[bmddm]["kmdy3"] = append(bkdMapping[bmddm]["kmdy3"], ksbh)
		}
		if strings.HasPrefix(ywk2m, "40") && wgym != "-" {
			bkdMapping[bmddm]["kmdy4"] = append(bkdMapping[bmddm]["kmdy4"], ksbh)
		}
	}
	fmt.Println(len(bkdMapping["2105"]["kmdy2"]), len(bkdMapping["2105"]["kmdy3"]), len(bkdMapping["2105"]["kmdy4"]))
	fmt.Println(len(bkdMapping["2106"]["kmdy2"]), len(bkdMapping["2106"]["kmdy3"]), len(bkdMapping["2106"]["kmdy4"]))
	fmt.Println(len(bkdMapping["2107"]["kmdy2"]), len(bkdMapping["2107"]["kmdy3"]), len(bkdMapping["2107"]["kmdy4"]))

	// 返回 DBF 文件的读取内容
	common.ResOk(c, "ok", nil)
}

func UploadBkdDBF(c *gin.Context) {
	year, err := common.GetCurrentYear()
	if err != nil {
		common.ResError(c, "获取年份信息失败")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		common.ResError(c, "获取文件失败")
		return
	}

	// 将文件保存到临时目录
	tempFile, err := ioutil.TempFile("", "uploaded-bkd-*.dbf")
	if err != nil {
		common.ResError(c, "创建文件失败")
		return
	}
	defer os.Remove(tempFile.Name()) // 处理完成后删除临时文件

	if err := c.SaveUploadedFile(file, tempFile.Name()); err != nil {
		common.ResError(c, "保存文件失败")
		return
	}

	// 打开 DBF 文件并使用 GBK 解码器
	dbfTable, err := dbf.OpenFile(tempFile.Name(), &GBKDecoder{})
	if err != nil {
		common.ResError(c, "读取文件失败")
		return
	}
	defer dbfTable.Close()

	// 遍历并输出记录内容
	var insertData []*model.CheckData
	var insertCheckData []*model.StepCheckData
	for !dbfTable.EOF() {
		record, err := dbfTable.Record()
		if err != nil {
			common.ResError(c, "读取记录失败")
			return
		}
		dbfTable.Skip(1)

		// 跳过已删除的记录
		if record.Deleted {
			continue
		}
		bmddm := strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("BMDDM"))).(string), " ", "")
		bmdmc := strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("BMDMC"))).(string), " ", "")
		insertData = append(insertData, &model.CheckData{
			Bmddm: bmddm,
			Bmdmc: bmdmc,
			Step:  model.CHECK_STEP_WAITING,
			Year:  year,
		})
		insertCheckData = append(insertCheckData, &model.StepCheckData{
			Bmddm:  bmddm,
			Step:   model.CHECK_STEP_WAITING,
			Status: model.CHECK_STATUS_WAITING,
			Year:   year,
		})
		if len(insertData) >= 1000 {
			_, err := conf.Mysql.Insert(&insertData)
			if err != nil {
				common.ResError(c, "写入数据失败")
				continue
			}
			_, err = conf.Mysql.Insert(&insertCheckData)
			if err != nil {
				common.ResError(c, "写入步骤数据失败")
				continue
			}
			insertData = []*model.CheckData{}
			insertCheckData = []*model.StepCheckData{}
		}
	}
	if len(insertData) > 0 {
		_, err := conf.Mysql.Insert(&insertData)
		if err != nil {
			common.ResError(c, "写入数据失败")
			return
		}
		_, err = conf.Mysql.Insert(&insertCheckData)
		if err != nil {
			common.ResError(c, "写入步骤数据失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}

func GetCheckList(c *gin.Context) {
	bmddm := c.Query("bmddm")
	bmdmc := c.Query("bmdmc")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.StepListSer(c, bmddm, bmdmc, pageInt, pageSizeInt)
}

func GetCheckInfo(c *gin.Context) {
	step := c.Query("step")
	stepInt, _ := strconv.Atoi(step)
	service.StepInfoSer(c, stepInt)
}

func Check(c *gin.Context) {
	var checkReq model.CheckReq
	if err := c.ShouldBindJSON(&checkReq); err != nil {
		log.Fatal(err)
		return
	}
	service.CheckSer(c, checkReq)
}

func NextStep(c *gin.Context) {
	var nextStepReq model.NextStepReq
	if err := c.ShouldBindJSON(&nextStepReq); err != nil {
		log.Fatal(err)
		return
	}
	service.NextSer(c, nextStepReq)
}
