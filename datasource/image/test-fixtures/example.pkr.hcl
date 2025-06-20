variable "access_key" {
  type    = string
  default = env("VOLCENGINE_ACCESS_KEY")
}

variable "secret_key" {
  type    = string
  default = env("VOLCENGINE_SECRET_KEY")
}

variable "region_id" {
  type    = string
  default = "cn-beijing"
}

variable "image_name" {
  type    = string
  default = "Ubuntu 20.04 with GPU Driver 570.86.15 and doca 2.5.0 64 bit"
}

data "volcengine-images" "example" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region_id
  tag_filters {
    key    = "test-packer"
    values = ["test-packer-value"]
  }
  image_name = var.image_name
}

source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = ["source.null.basic-example"]

  provisioner "shell-local" {
    inline = [
      "echo image_id:            ${data.volcengine-images.example.images[0].image_id}",
      "echo image_name:          ${data.volcengine-images.example.images[0].image_name}",
      "echo description:         ${data.volcengine-images.example.images[0].description}",
      "echo platform:            ${data.volcengine-images.example.images[0].platform}",
      "echo platform_version:    ${data.volcengine-images.example.images[0].platform_version}",
      "echo visibility:          ${data.volcengine-images.example.images[0].visibility}",
      "echo is_support_cloud_init: ${data.volcengine-images.example.images[0].is_support_cloud_init}"
    ]
  }
}
