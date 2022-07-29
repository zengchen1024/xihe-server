#!/bin/bash

# set -euo pipefail

work_dir=$(pwd)
target_dir=$work_dir/{TARGET}

gitlab_endpoint=${GITLAB_ENDPOINT-""}
xihe_user=${XIHE_USER-""}
xihe_user_token=${XIHE_USER_TOKEN-""}

f="$target_dir/config.json"
pretain_file=""

if [ -e "$f" -a -s "$f" ]; then
	v=$(python3 ./pretrain.py $f)
	if [ $? -ne 0 ]; then
		echo $v
		exit 1
	fi

	if [ -n "$v" ]; then
		owner=$(echo $v | sed -n '1p')
		repo=$(echo $v | sed -n '2p')
		pretain_file=$(echo $v | sed -n '3p')

		mkdir $target_dir/$owner
		cd $target_dir/$owner
		git clone http://${xihe_user}:${xihe_user_token}@${gitlab_endpoint}/${owner}/${repo}
	fi
fi

cd $target_dir

python3 ./app.py
