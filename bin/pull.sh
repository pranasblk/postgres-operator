#!/bin/bash 

# Copyright 2017-2018 Crunchy Data Solutions, Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

export P_IMAGE_PREFIX=registry-dev.crunchydata.com/crunchydata
export CO_IMAGE_TAG=centos7-2.6

docker pull $P_IMAGE_PREFIX/pgo-lspvc:$CO_IMAGE_TAG
docker tag $P_IMAGE_PREFIX/pgo-lspvc:$CO_IMAGE_TAG $CO_IMAGE_PREFIX/pgo-lspvc:$CO_IMAGE_TAG
docker pull $P_IMAGE_PREFIX/postgres-operator:$CO_IMAGE_TAG   
docker tag $P_IMAGE_PREFIX/postgres-operator:$CO_IMAGE_TAG $CO_IMAGE_PREFIX/postgres-operator:$CO_IMAGE_TAG   
docker pull $P_IMAGE_PREFIX/pgo-load:$CO_IMAGE_TAG   
docker tag $P_IMAGE_PREFIX/pgo-load:$CO_IMAGE_TAG $CO_IMAGE_PREFIX/pgo-load:$CO_IMAGE_TAG   
docker pull $P_IMAGE_PREFIX/pgo-apiserver:$CO_IMAGE_TAG  
docker tag $P_IMAGE_PREFIX/pgo-apiserver:$CO_IMAGE_TAG $CO_IMAGE_PREFIX/pgo-apiserver:$CO_IMAGE_TAG  
docker pull $P_IMAGE_PREFIX/pgo-rmdata:$CO_IMAGE_TAG  
docker tag $P_IMAGE_PREFIX/pgo-rmdata:$CO_IMAGE_TAG $CO_IMAGE_PREFIX/pgo-rmdata:$CO_IMAGE_TAG  
