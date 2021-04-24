#!/bin/bash

docker run -p 27017:27017 -it -v mongodata:/data/db --name mongodb -d mongo