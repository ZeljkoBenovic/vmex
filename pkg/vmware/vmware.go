package vmware

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
)

type Vmware struct {
	mng *view.Manager
	vm  *view.ContainerView
	vms []mo.VirtualMachine
}

type VMData struct {
	Data          []VM
	TotalMemoryGB float64
	TotalHddGB    float64
	TotalCPU      int32
}

type VM struct {
	Name       string
	Hostname   string
	TotalHDD   string
	UsedHDD    string
	Memory     string
	CPU        int
	IPAddress  string
	UptimeHour string
}

type Options struct {
	Filter string
}

type Opts func(*Options)

func New(host, user, pass string) (*Vmware, error) {
	vcUrl, err := soap.ParseURL(host)
	if err != nil {
		return nil, fmt.Errorf("could not parse host: %w", err)
	}

	vcUrl.User = url.UserPassword(user, pass)

	sess := cache.Session{
		URL:      vcUrl,
		Insecure: true,
	}

	cl := new(vim25.Client)
	if err = sess.Login(context.Background(), cl, nil); err != nil {
		return nil, fmt.Errorf("could not login: %w", err)
	}

	v := &Vmware{}
	v.mng = view.NewManager(cl)

	vm, err := v.mng.CreateContainerView(context.Background(), cl.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return nil, fmt.Errorf("could not create virtual machine view: %w", err)
	}

	v.vm = vm

	if err = v.vm.Retrieve(context.Background(), []string{"VirtualMachine"}, []string{"summary"}, &v.vms); err != nil {
		return nil, fmt.Errorf("could not retrieve the list of vms: %w", err)
	}

	return v, nil
}

func (v *Vmware) GetAll(opts ...Opts) VMData {
	var (
		vmData = VMData{}
		opt    = &Options{Filter: ""}
	)

	for _, f := range opts {
		f(opt)
	}

	filters := strings.Split(opt.Filter, ",")

	for _, filter := range filters {
		for _, vm := range v.vms {
			if strings.Contains(strings.ToUpper(vm.Summary.Config.Name), strings.ToUpper(strings.TrimSpace(filter))) {
				vmName := vm.Summary.Config.Name
				provisionedHddGB := float64(vm.Summary.Storage.Committed+vm.Summary.Storage.Uncommitted) / 1024 / 1024 / 1024
				usedHddGB := float64(vm.Summary.Storage.Uncommitted) / 1024 / 1024 / 1024
				ramGB := float64(vm.Summary.Config.MemorySizeMB) / 1024
				uptimeHour := float64(vm.Summary.QuickStats.UptimeSeconds / 60 / 60)
				cpuNum := vm.Summary.Config.NumCpu
				ipAddress := vm.Summary.Guest.IpAddress
				hostName := vm.Summary.Guest.HostName

				vmd := VM{
					Name:       vmName,
					Hostname:   hostName,
					TotalHDD:   fmt.Sprintf("%.2fGB", provisionedHddGB),
					UsedHDD:    fmt.Sprintf("%.2fGB", usedHddGB),
					Memory:     fmt.Sprintf("%.2fGB", ramGB),
					CPU:        int(cpuNum),
					IPAddress:  fmt.Sprintf("%s", ipAddress),
					UptimeHour: fmt.Sprintf("%.2fh", uptimeHour),
				}

				vmData.TotalHddGB += provisionedHddGB
				vmData.TotalMemoryGB += ramGB
				vmData.TotalCPU += cpuNum

				vmData.Data = append(vmData.Data, vmd)
			}
		}
	}

	return vmData
}

func WithFilter(filterString string) Opts {
	return func(options *Options) {
		options.Filter = filterString
	}
}
