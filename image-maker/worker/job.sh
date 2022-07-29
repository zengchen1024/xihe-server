#!/bin/sh

work_dir=$(pwd)
gradio_dir=${work_dir}/gradio

gitlab_endpoint=${GITLAB_ENDPOINT-""}
xihe_user=${XIHE_USER-""}
xihe_user_token=${XIHE_USER_TOKEN-""}
repo=${XIHE_REPO-""}
commit=${XIHE_REPO_COMMIT-""}

if [ -z "$gitlab_endpoint" ]; then
	echo "no gitlab endpoint"
	exit 1
fi

if [ -z "$xihe_user" ]; then
	echo "no xihe uer"
	exit 1
fi

if [ -z "$xihe_user_token" ]; then
	echo "no xihe uer token"
	exit 1
fi

if [ -z "$repo" ]; then
	echo "no xihe repo"
	exit 1
fi

if [ -z "$commit" ]; then
	echo "no xihe commit"
	exit 1
fi

# clone repo
echo "start clone ${xihe_user}/${repo}"

git clone http://${xihe_user}:${xihe_user_token}@${gitlab_endpoint}/${xihe_user}/${repo}

echo "end clone"

cd $repo
commit_id=$(git describe --tags --always --dirty)
cd $work_dir

if [ "$commit" != "$commit_id" ]; then
	echo "it is not the expect commit, and the repo has changed."
	exit 1
fi

# prepare base image
echo "start podman pull"

mount --make-rshared /

podman pull docker.io/library/python:3.9.13

echo "end podman pull"

# make image
mv $repo $gradio_dir

cd $gradio_dir

sed -i "s/{TARGET}/""$repo/" run.sh
sed -i "s/{TARGET}/""$repo/g" Dockerfile

msg=$(podman build -t $repo:$commit -f ./Dockerfile)
if [ $? -ne 0 ]; then
	echo $msg
	echo "failed"
	exit 1
fi

cd $work_dir

# test image

echo "start test image"

cat << EOF >> env.list
XIHE_USER=$XIHE_USER
XIHE_USER_TOKEN=$XIHE_USER_TOKEN
GITLAB_ENDPOINT=$GITLAB_ENDPOINT
EOF

podman run --env-file=./env.list $repo:$commit

echo "end test image"

# send result

echo "success"
