package template

var RunTemplate = "#ps -ef|grep %s|grep -v grep|awk '{print $2}'|xargs kill\npid=`ps -ef|grep %s|grep -v grep|awk '{print $2}'`\nif [ -n \"$pid\" ]; then\n    echo \"process id is: $pid\"\n    kill $pid\n    echo 'process killed'\nelse\n    echo \"%s process id is not found!\"\nfi\nexport GIN_MODE=release\necho 'starting...'\n#nohup ./%s &\npm2 start pm2.yml\necho 'started!'"
