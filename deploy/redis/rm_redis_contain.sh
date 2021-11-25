docker rm -f $(docker ps -a | grep redis | awk '{print $1}')
