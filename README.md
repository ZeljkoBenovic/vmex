# VMEX
A small tool used to export VM list from `VMWAREvCenter` into the Excel spreadsheet.     
The export functionality in the `vCenter` itself, exports all VMs into a `csv` file, which needs to be processed.   
This tool aims to ease this export procedure, by exporting (`all` or `filtered by name`) VM into a nicely formatted Excel file.    

## Usage
### Flags
* `host` - `vCenter` host url
* `user` - `vCenter` username
* `pass` - `vCenter` password
* `filter` (optional) - comma delimited string which will be used to filter out VM names
* `path` (optional) - the path and filename where the report will be saved

### Example
```bash
vmex -host "https://vcenter.host.local" -user "administrator@vcenter.local" -pass "xxxxxxxx" -filter "web" -path web_server_vms.xlsx
```