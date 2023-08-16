package template

var ReadMeTemplate = `# %s

config: stores configuration functions
common: stores helper functions
docs: stores Swagger documentation
codes: stores error codes
enums: stores constant codes
hooks: stores asynchronous hooks
middleware: stores middleware
controller: stores controllers
model: stores database models
pkg: stores third-party packages
repository: stores database interaction layer
requests: stores request models
responses: stores response models
routers: stores routers
sysinit: stores system initialization
services: stores service layer
jobs: stores timer jobs
`
