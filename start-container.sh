#!/bin/sh
echo "Enter container id: "  
read cid  
sudo docker run -p 8888:5000 $cid
