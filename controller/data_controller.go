package controller

import (
	"bytes"
	"data_verify/common"
	"data_verify/conf"
	"data_verify/model"
	dbf "github.com/SebastiaanKlippert/go-foxpro-dbf"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
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

// UploadAndReadDBF 上传文件并读取内容
func UploadAndReadDBF(c *gin.Context) {
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

	// 遍历并输出记录内容
	var insertData []*model.SbkData
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

		// 根据字段名读取特定字段
		insertData = append(insertData, &model.SbkData{
			Idxx:    strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("IDXX"))).(string), " ", ""),
			Bxyj:    strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("BXYJ"))).(string), " ", ""),
			Bmddm:   strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("BMDDM"))).(string), " ", ""),
			Bmdmc:   strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("BMDMC"))).(string), " ", ""),
			Hbbmdmc: strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("HBBMDMC"))).(string), " ", ""),
			Xm:      strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("XM"))).(string), " ", ""),
			Ksbh:    strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("KSBH"))).(string), " ", ""),
			Kmdm:    strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("KMDM"))).(string), " ", ""),
			Kmmc:    strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("KMMC"))).(string), " ", ""),
			Kmdy:    strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("KMDY"))).(string), " ", ""),
			Dy2:     common.IgnoreError(record.Field(dbfTable.FieldPos("DY2"))).(int64),
			Dy3:     common.IgnoreError(record.Field(dbfTable.FieldPos("DY3"))).(int64),
			Dy4:     common.IgnoreError(record.Field(dbfTable.FieldPos("DY4"))).(int64),
			Kssj:    strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("KSSJ"))).(string), " ", ""),
			Xlsh:    common.IgnoreError(record.Field(dbfTable.FieldPos("XLSH"))).(int64),
			Dlsh:    common.IgnoreError(record.Field(dbfTable.FieldPos("DLSH"))).(int64),
			Djtxm1:  strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("DJTXM1"))).(string), " ", ""),
			Djtxm2:  strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("DJTXM2"))).(string), " ", ""),
			Djtxm3:  strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("DJTXM3"))).(string), " ", ""),
			Djtxm4:  strings.ReplaceAll(common.IgnoreError(record.Field(dbfTable.FieldPos("DJTXM4"))).(string), " ", ""),
			Year:    2024,
		})

		if len(insertData) >= 1000 {
			_, err := conf.Mysql.Insert(&insertData)
			if err != nil {
				common.ResError(c, "写入数据失败")
				continue
			}
			insertData = []*model.SbkData{}
		}
	}
	if len(insertData) > 0 {
		_, err := conf.Mysql.Insert(&insertData)
		if err != nil {
			common.ResError(c, "写入数据失败")
			return
		}
	}

	// 返回 DBF 文件的读取内容
	common.ResOk(c, "ok", nil)
}
