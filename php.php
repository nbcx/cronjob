<?php
sleep(3);
echo "hello";
error_log(date('Y-m-d H:i:s')." [debug] ".time()."\n", 3, "/Users/collin/Server/www/ol.cx/phpMobi/log.txt");
echo "接收到{$argc}个参数\n";
print_r($argv);