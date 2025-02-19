package controller

import (
	"data_verify/common"
	"fmt"
	dbf "github.com/SebastiaanKlippert/go-foxpro-dbf"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx/v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CompareData(c *gin.Context) {
	//year, err := common.GetCurrentYear()
	//if err != nil {
	//	common.ResError(c, "获取年份信息失败")
	//	return
	//}
	// 获取上传的文件
	form, err := c.MultipartForm()
	if err != nil {
		common.ResError(c, "获取文件失败")
		return
	}

	files := form.File["files"]

	file1 := files[0]
	file2 := files[1]

	ext1 := strings.TrimPrefix(filepath.Ext(file1.Filename), ".")

	// 将文件保存到临时目录
	tempFile1, err := ioutil.TempFile("", "uploaded-*."+ext1)
	if err != nil {
		common.ResError(c, "创建文件失败")
		return
	}
	defer os.Remove(tempFile1.Name()) // 处理完成后删除临时文件

	if err := c.SaveUploadedFile(file1, tempFile1.Name()); err != nil {
		common.ResError(c, "保存文件失败")
		return
	}

	ext2 := strings.TrimPrefix(filepath.Ext(file2.Filename), ".")

	// 将文件保存到临时目录
	tempFile2, err := ioutil.TempFile("", "uploaded-*."+ext2)
	if err != nil {
		common.ResError(c, "创建文件失败")
		return
	}
	defer os.Remove(tempFile2.Name()) // 处理完成后删除临时文件

	if err := c.SaveUploadedFile(file2, tempFile2.Name()); err != nil {
		common.ResError(c, "保存文件失败")
		return
	}

	var file1Data [][]interface{}
	var file2Data [][]interface{}

	if ext1 == "dbf" {
		dbfTable, err := dbf.OpenFile(tempFile1.Name(), &GBKDecoder{})
		if err != nil {
			common.ResError(c, "读取文件失败")
			return
		}
		defer dbfTable.Close()
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
			var insertData []interface{}
			for _, d := range record.FieldSlice() {
				insertData = append(insertData, d)
			}
			file1Data = append(file1Data, insertData)
			insertData = []interface{}{}
		}
	}

	if ext1 == "xlsx" {
		xlsxTable, err := xlsx.OpenFile(tempFile1.Name())
		if err != nil {
			common.ResError(c, "读取文件失败")
			return
		}
		for _, sheet := range xlsxTable.Sheets {
			// 遍历工作表中的每一行
			err := sheet.ForEachRow(func(row *xlsx.Row) error {
				var insertData []interface{}

				// 遍历行中的每一列
				err := row.ForEachCell(func(cell *xlsx.Cell) error {
					// 获取单元格的值
					value, err := cell.FormattedValue()
					if err != nil {
						return err
					}
					insertData = append(insertData, value)
					return nil
				})
				if err != nil {
					return err
				}

				// 将当前行数据添加到 file1Data
				file1Data = append(file1Data, insertData)
				return nil
			})
			if err != nil {
				log.Fatalf("读取数据失败")
				return
			}
		}
	}

	if ext2 == "dbf" {
		dbfTable, err := dbf.OpenFile(tempFile2.Name(), &GBKDecoder{})
		if err != nil {
			common.ResError(c, "读取文件失败")
			return
		}
		defer dbfTable.Close()
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
			var insertData []interface{}
			for _, d := range record.FieldSlice() {
				insertData = append(insertData, d)
			}
			file2Data = append(file2Data, insertData)
			insertData = []interface{}{}
		}
	}

	if ext2 == "xlsx" {
		xlsxTable, err := xlsx.OpenFile(tempFile2.Name())
		if err != nil {
			common.ResError(c, "读取文件失败")
			return
		}
		for _, sheet := range xlsxTable.Sheets {
			// 遍历工作表中的每一行
			err := sheet.ForEachRow(func(row *xlsx.Row) error {
				var insertData []interface{}

				// 遍历行中的每一列
				err := row.ForEachCell(func(cell *xlsx.Cell) error {
					// 获取单元格的值
					value, err := cell.FormattedValue()
					if err != nil {
						return err
					}
					insertData = append(insertData, value)
					return nil
				})
				if err != nil {
					return err
				}

				// 将当前行数据添加到 file1Data
				file2Data = append(file2Data, insertData)
				return nil
			})
			if err != nil {
				log.Fatalf("读取数据失败")
				return
			}
		}
	}

	if len(file1Data) != len(file2Data) {
		common.ResForbidden(c, "文件行数不同")
		return
	}

	if len(file1Data[0]) != len(file2Data[0]) {
		common.ResForbidden(c, "文件列数不同")
		return
	}

	type resObj struct {
		ErrMsg    string `json:"err_msg"`
		FirstErr  string `json:"first_err"`
		SecondErr string `json:"second_err"`
	}
	x := 0
	y := 0
	var errData []resObj
	for x < len(file1Data) {
		for y < len(file1Data[x]) {
			if file1Data[x][y] != file2Data[x][y] {
				errData = append(errData, resObj{
					ErrMsg:    fmt.Sprintf("%d行%d列出现错误！", x+1, y+1),
					FirstErr:  fmt.Sprintf("%s", file1Data[x][y]),
					SecondErr: fmt.Sprintf("%s", file2Data[x][y]),
				})
			}
			y++
		}
		x++
		y = 0
	}
	common.ResOk(c, "ok", errData)
}
