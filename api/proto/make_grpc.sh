#!/bin/bash
cmd="protoc " 
for d in $(ls -l | grep ^d | grep -v tests | awk '{print $9}')
do    
    cmd=$cmd" -I"$d
done
echo "proto dest dir : " $1
cmd=$cmd" --go_out=plugins=grpc:"$1" "
for f in $(find . -name "*.proto" -exec basename {} \;)
do
    if [ "$f" = "$http" -o "$f" = "$annotations" -o "$f" = "$descriptor" ];then
        continue
    else
        cmd=$cmd" "$f
    fi
    #echo $f
done
echo $cmd
$cmd

