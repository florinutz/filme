variable "ver" {}

variable "bucket" {
  description = "bucket to upload the build to"
  default = "lambda-example-flo"
}

variable "region" {
  default = "eu-central-1"
}

variable "file" {
  default = "build.zip"
}

variable "func_name" {
  default = "ServerlessExample"
}
