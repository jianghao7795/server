package example

import exampleServer "server/service/example"

type FileUploadAndDownloadApi struct{}
type ExcelApi struct{}

var excelService = exampleServer.ExcelServiceApp

type CustomerApi struct{}

var customerService = exampleServer.CustomerServiceApp
var fileUploadAndDownloadService = exampleServer.FileUploadAndDownloadServiceApp
