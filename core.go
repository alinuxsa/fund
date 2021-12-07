package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Data struct {
	Datas []struct {
		Fcode           string      `json:"FCODE"`
		Shortname       string      `json:"SHORTNAME"`
		Pdate           string      `json:"PDATE"`
		Nav             string      `json:"NAV"`
		Accnav          string      `json:"ACCNAV"`
		Navchgrt        string      `json:"NAVCHGRT"`
		Gsz             string      `json:"GSZ"`
		Gszzl           string      `json:"GSZZL"`
		Gztime          string      `json:"GZTIME"`
		Newprice        interface{} `json:"NEWPRICE"`
		Changeratio     interface{} `json:"CHANGERATIO"`
		Zjl             interface{} `json:"ZJL"`
		Hqdate          interface{} `json:"HQDATE"`
		Ishaveredpacket bool        `json:"ISHAVEREDPACKET"`
	} `json:"Datas"`
	ErrCode      int         `json:"ErrCode"`
	Success      bool        `json:"Success"`
	ErrMsg       interface{} `json:"ErrMsg"`
	Message      interface{} `json:"Message"`
	ErrorCode    string      `json:"ErrorCode"`
	ErrorMessage interface{} `json:"ErrorMessage"`
	ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
	TotalCount   int         `json:"TotalCount"`
	Expansion    struct {
		Gztime string `json:"GZTIME"`
		Fsrq   string `json:"FSRQ"`
	} `json:"Expansion"`
}

func GetFunds(funds []string) {
	codes := strings.Join(funds, ",")
	baseUrl := "https://fundmobapi.eastmoney.com/FundMNewApi/FundMNFInfo?pageIndex=1&pageSize=50&plat=Android&appType=ttjj&product=EFund&Version=1&deviceid=166b2763-375f-4af8-99cb-c12f08e0b88d&Fcodes=%s"
	url := fmt.Sprintf(baseUrl, codes)
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36 Edg/89.0.774.68")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	ret := new(Data)
	err = json.Unmarshal(data, ret)
	if err != nil {
		panic(err)
	}

	updateTime := ret.Datas[0].Gztime

	tableData := [][]string{}
	for _, f := range ret.Datas {
		lst := []string{}
		lst = append(lst, f.Fcode, f.Shortname, f.Nav, f.Gszzl)
		tableData = append(tableData, lst)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"基金代码", "基金名字", "净值", "估算涨幅"})
	table.SetFooter([]string{"", "", "更新时间", updateTime})
	table.SetBorder(true)

	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.Bold})

	colorUp := []tablewriter.Colors{{}, {}, {}, {tablewriter.Bold, tablewriter.FgRedColor}}
	colorDown := []tablewriter.Colors{{}, {}, {}, {tablewriter.Bold, tablewriter.FgGreenColor}}

	for _, row := range tableData {
		if strings.HasPrefix(row[3], "-") {
			table.Rich(row, colorDown)
		} else {
			table.Rich(row, colorUp)
		}

	}
	table.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_RIGHT})
	table.Render()

}
