package excel

import (
	"fmt"

	"vmware-api/pkg/vmware"

	"github.com/xuri/excelize/v2"
)

type Excel struct {
	file   *excelize.File
	vmData vmware.VMData
}

type Options struct {
	FilePath string
}

type Opts func(*Options)

func New(vmData vmware.VMData) *Excel {
	return &Excel{
		file:   excelize.NewFile(),
		vmData: vmData,
	}
}

func (e *Excel) CreateTable(opts ...Opts) error {
	var (
		o             = &Options{FilePath: "vm_list.xlsx"}
		columnTracker = 0
	)
	for _, f := range opts {
		f(o)
	}

	e.file.SetSheetName("Sheet1", "VM LIST")

	e.file.SetCellValue("VM LIST", "A1", "NAME")
	e.file.SetCellValue("VM LIST", "B1", "HOSTNAME")
	e.file.SetCellValue("VM LIST", "C1", "HDD")
	e.file.SetCellValue("VM LIST", "D1", "RAM")
	e.file.SetCellValue("VM LIST", "E1", "CPU")
	e.file.SetCellValue("VM LIST", "F1", "IP")
	//e.file.SetCellValue("VM LIST", "G1", "UPTIME")

	e.file.SetColWidth("VM LIST", "A", "B", 30)
	e.file.SetColWidth("VM LIST", "C", "C", 12)
	e.file.SetColWidth("VM LIST", "D", "E", 8)

	style, _ := e.file.NewStyle(&excelize.Style{
		Border: nil,
		Fill:   excelize.Fill{},
		Font:   nil,
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	e.file.SetCellStyle("VM LIST", "A1", "F1", style)

	columnTracker = 2

	for i, vm := range e.vmData.Data {
		e.file.SetCellValue("VM LIST", fmt.Sprintf("A%d", i+2), vm.Name)
		e.file.SetCellValue("VM LIST", fmt.Sprintf("B%d", i+2), vm.Hostname)
		e.file.SetCellValue("VM LIST", fmt.Sprintf("C%d", i+2), vm.TotalHDD)
		e.file.SetCellValue("VM LIST", fmt.Sprintf("D%d", i+2), vm.Memory)
		e.file.SetCellInt("VM LIST", fmt.Sprintf("E%d", i+2), vm.CPU)
		e.file.SetCellValue("VM LIST", fmt.Sprintf("F%d", i+2), vm.IPAddress)

		e.file.SetCellStyle("VM LIST", fmt.Sprintf("E%d", i+2), fmt.Sprintf("E%d", i+2), style)

		columnTracker++
	}

	columnTracker += 1

	e.file.SetCellValue("VM LIST", fmt.Sprintf("A%d", columnTracker), "TOTAL")
	e.file.SetCellValue("VM LIST", fmt.Sprintf("C%d", columnTracker), e.vmData.TotalHddGB)
	e.file.SetCellValue("VM LIST", fmt.Sprintf("D%d", columnTracker), e.vmData.TotalMemoryGB)
	e.file.SetCellValue("VM LIST", fmt.Sprintf("E%d", columnTracker), e.vmData.TotalCPU)

	e.file.SetCellStyle("VM LIST", fmt.Sprintf("D%d", columnTracker), fmt.Sprintf("E%d", columnTracker), style)

	e.file.AddTable("VM LIST", &excelize.Table{
		Range:     fmt.Sprintf("A1:F%d", columnTracker),
		Name:      "table",
		StyleName: "TableStyleMedium7",
	})

	if err := e.file.SaveAs(o.FilePath); err != nil {
		return fmt.Errorf("could not save report: %w", err)
	}

	fmt.Printf("Report saved at: %s", o.FilePath)

	return nil
}

func WithFilePath(filePath string) Opts {
	return func(options *Options) {
		options.FilePath = filePath
	}
}
