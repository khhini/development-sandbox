#/bin/bash
export BRANCH_NAME=main
export LOCATION=asia-east1
export SHORT_SHA=5d74f0a
export _DEPLOYMENT_ENV=dev
export _BUILD_PATH=./golang/gotodo-api
export _DEPLOYMENT_CONFIG_PATH=${_BUILD_PATH}/deployments
export _APP_NAME=gotodo
export _APP_SERVICE=api
export _CLOUD_RUN_SERVICE_YAML_PATH=${_DEPLOYMENT_CONFIG_PATH}/cloudrun/services.yaml
export _CLOUD_RUN_SERVICE_PROJECT=khhini-development-sandbox
export _ARTIFACT_REGISTRY_URI=asia-east1-docker.pkg.dev/khhini-devops-2705/docker-repo
export _CLOUDBUILD_CD_TRIGGER="34f55dd2-1d74-4719-9db0-272624fa6ee6"

gcloud builds triggers run ${_CLOUDBUILD_CD_TRIGGER} \
  --substitutions=_APP_NAME=${_APP_NAME},_APP_SERVICE=${_APP_SERVICE},_DEPLOYMENT_ENV=${_DEPLOYMENT_ENV},_CLOUD_RUN_SERVICE_YAML_PATH=${_CLOUD_RUN_SERVICE_YAML_PATH},_ARTIFACT_REGISTRY_URI=${_ARTIFACT_REGISTRY_URI},_IMAGE_VERSION_TAG=${SHORT_SHA} \
  --branch=${BRANCH_NAME} \
  --region=${LOCATION}
