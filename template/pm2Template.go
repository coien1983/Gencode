package template

var Pm2Template = "apps:\n  - name: %s\n    instances: 1\n    exec_mode: fork\n    interpreter: \"./%s\"\n    interpreter_args: \"-f ./config/config-online.yaml\"\n    script: \"./pm2.yml\""
