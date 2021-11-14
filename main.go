package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/fdxxw/go-wen"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/cast"
	"github.com/thoas/go-funk"
)

func main() {
	flag.Parse()
	assetIds := flag.Args()
	if len(assetIds) == 0 {
		assetIds = []string{"86", "102", "79", "58", "84", "78"}
	}
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table := widgets.NewTable()

	table.Title = "Table"
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.SetRect(0, 0, 60, 20)
	table.RowSeparator = true
	table.TextAlignment = ui.AlignCenter

	update := func() {
		otcs, _ := fetch()
		rows := [][]string{{"assetId", "sellPriceUsdt"}}
		for _, otc := range otcs {
			if funk.ContainsString(assetIds, cast.ToString(otc.AssetId)) {
				rows = append(rows, []string{cast.ToString(otc.AssetId), otc.SellPriceUsdt})
			}
		}
		table.Rows = rows
		ui.Render(table)
	}

	update()
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(5 * time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			update()
		}
	}
}

func fetch() ([]Otc, error) {
	url := "https://app.eiduwejdk.com/api/v1/otc"
	resp, err := wen.HttpGet(url, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var typ struct {
		Code int `json:"code"`
		Data []Otc
	}
	json.Unmarshal(data, &typ)
	return typ.Data, nil
}

type Otc struct {
	ID            int     `json:"id"`
	AssetId       int     `json:"assetId"`
	AssetUuid     string  `json:"assetUuid"`
	BuyPrice      float64 `json:"buyPrice"`
	SellPrice     float64 `json:"sellPrice"`
	BuyPriceUsdt  float64 `json:"buyPriceUsdt"`
	SellPriceUsdt string  `json:"sellPriceUsdt"`
	UsdtByPrice   float64 `json:"usdtByPrice"`
	UsdtSellPrice float64 `json:"usdtSellPrice"`
	FloatingRate  float64 `json:"floatingRate"`
	Grade         float64 `json:"grade"`
	MinimumAmount float64 `json:"minimumAmount"`
}
