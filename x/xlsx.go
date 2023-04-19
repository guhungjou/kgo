package x

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx/v3"
)

func ParseXLSXSheet(sh *xlsx.Sheet, value interface{}) error {
	if sh.MaxRow <= 1 {
		return nil
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr || v.IsNil() || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("invalid ptr")
	}

	titleMap := make(map[string]int)
	for i := 0; i < sh.MaxCol; i++ {
		cell, _ := sh.Cell(0, i)
		title := cell.String()
		if title == "" {
			continue
		}
		titleMap[title] = i
	}

	typ := v.Elem().Type()
	e := typ.Elem().Elem()

	vvv := v.Elem()

	for i := 1; i < sh.MaxRow; i++ {
		vvvv := reflect.New(e)
		vv := vvvv.Elem()
		inserted := false
		// if cell, err := sh.Cell(i, 0); err != nil {
		// 	break
		// } else if cell.String() == "" {
		// 	break
		// }
		for j := 0; j < vv.NumField(); j++ {
			tfield := vv.Type().Field(j)
			vfield := vv.Field(j)
			if xlsx := tfield.Tag.Get("xlsx"); xlsx != "" {
				seps := strings.Split(xlsx, ",")
				title := xlsx
				format := "2006/1/2 15:04:05"
				if len(seps) > 1 {
					title = seps[0]
					format = seps[1]
				}
				if col, ok := titleMap[title]; ok {
					cell, _ := sh.Cell(i, col)
					if cell != nil {
						str := strings.TrimSpace(cell.String())
						if str != "" {
							switch vfield.Kind() {
							case reflect.String:
								vfield.SetString(str)
							case reflect.Int64, reflect.Int, reflect.Int32:
								iv, _ := strconv.ParseInt(str, 10, 64)
								vfield.SetInt(iv)
							case reflect.TypeOf(time.Time{}).Kind():
								// t, _ := time.ParseInLocation(format, str, time.Local)
								t, _ := cell.GetTime(false)
								if t.IsZero() {
									t, _ = time.ParseInLocation(format, str, time.Local)
								}
								vfield.Set(reflect.ValueOf(t))
							case reflect.Float64:
								iv, _ := strconv.ParseFloat(str, 64)
								vfield.SetFloat(iv)
							default:
								continue
							}
							inserted = true
						}
					}
				}
			}
		}
		if inserted {
			vvv = reflect.Append(vvv, vvvv)
		}
	}

	v.Elem().Set(vvv)

	return nil
}

func CreateXLSXBook(sheetName string, headers []string, rows [][]interface{}) (*xlsx.File, error) {
	wb := xlsx.NewFile()
	sh, err := wb.AddSheet(sheetName)
	if err != nil {
		return nil, err
	}
	sh.SetColWidth(0, 100, 15)
	sh.SheetFormat.DefaultRowHeight = 18

	row := sh.AddRow()
	headStyle := xlsx.NewStyle()
	headStyle.Alignment.Horizontal = "center"
	headStyle.Font.Color = "#00CCFF00"
	headStyle.Font.Size = 12
	headStyle.Font.Bold = true
	headStyle.ApplyAlignment = true
	headStyle.ApplyFont = true
	headStyle.ApplyFill = false
	for _, header := range headers {
		cell := row.AddCell()
		cell.SetString(header)
		cell.SetStyle(headStyle)
	}

	for _, row := range rows {
		r := sh.AddRow()
		for _, col := range row {
			switch v := col.(type) {
			case string:
				r.AddCell().SetString(v)
			case time.Time:
				if !v.IsZero() {
					if v.Hour() == 0 && v.Minute() == 0 && v.Second() == 0 {
						r.AddCell().SetDateWithOptions(v, xlsx.DateTimeOptions{Location: time.Local, ExcelTimeFormat: "yyyy-mm-dd"})
					} else {
						r.AddCell().SetDateWithOptions(v, xlsx.DateTimeOptions{Location: time.Local, ExcelTimeFormat: "yyyy-mm-dd hh:MM:ss"})
					}
				}
			case int:
				r.AddCell().SetInt(v)
			case int64:
				r.AddCell().SetInt64(v)
			case float32:
				r.AddCell().SetFloat(float64(v))
			case float64:
				r.AddCell().SetFloat(v)
			default:
				return nil, fmt.Errorf("unknown cell type")
			}
		}
	}

	return wb, nil
}
