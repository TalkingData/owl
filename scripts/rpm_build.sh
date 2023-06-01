#!/bin/bash

pwd=`pwd`
src="$pwd/../src"
go_output="$pwd/../output"
rpmbuild_path="/root/rpmbuild"

branch=$(git symbolic-ref --short -q HEAD)
commit=$(git rev-parse --short HEAD)

service_name=$1
version=$2


function print_usage() {
  echo "USAGE:"
  echo $'\t'"$0 <ServiceName> <Version>"

  echo "e.g.:"
  echo $'\t'"$0 api 6.0.0"
  echo $'\t'"$0 agent 6.0.0"

  exit 1;
}

function build_go() {
  echo "[I] Prepare building: owl-$service_name, version: $version, output: $go_output"
  # 切换目录经并开始构建
  cd "$src/$service_name"
  CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-w\
    -X owl/common/global.Version=$version\
    -X owl/common/global.Branch=$branch\
    -X owl/common/global.Commit=$commit" -o "$go_output/owl-$service_name-$version"

  # 构建是否成功
  if [ $? -eq 0 ]; then
    echo "[I] Success build owl-$service_name"
  else
    echo "[E] Failed build owl-$service_name"
    exit 2
  fi
}

function build_rpm() {
  echo "[I] Prepare building rpm: owl-$service_name-$version.rpm"

  # 准备拷贝spec文件
  echo "[I] Copy spec file to rpmbuild path"
  spec_file_pathname="$rpmbuild_path/SPECS/owl-$service_name-$version.spec"
  # 如果spec文件存在，就删除
  if [ -f "$spec_file_pathname" ]; then
   rm -rf "$spec_file_pathname"
  fi

  # 拷贝spec文件
  cp "$pwd/specs/owl-$service_name.spec" "$spec_file_pathname"

  # 替换spec文件内版本
  sed -i "s/Version: .*/Version:   $version/g" "$spec_file_pathname"

  # 准备拷贝source文件（bin文件、conf文件和service文件）
  echo "[I] Copy bin, conf, service files to rpmbuild path"
  rpm_source_dir="$rpmbuild_path/SOURCES"
  service_rpm_source_dir="$rpm_source_dir/owl-$service_name-$version"
  # 如果source目录存在，就删除
  if [ -f "$service_rpm_source_dir" ]; then
   rm -rf "$service_rpm_source_dir"
  fi
  if [ -d "$service_rpm_source_dir" ]; then
   rm -rf "$service_rpm_source_dir"
  fi
  # 重新创建source目录
  mkdir -p "$service_rpm_source_dir/bin"
  mkdir -p "$service_rpm_source_dir/conf"

  # 拷贝bin文件
  cp "$go_output/owl-$service_name-$version" "$service_rpm_source_dir/bin/owl-$service_name"
  # 拷贝conf文件
  cp "$src/conf/owl_$service_name.conf-sample" "$service_rpm_source_dir/conf/owl_$service_name.conf"
  # 拷贝service文件
  cp "$pwd/systemd/owl-$service_name.service" "$service_rpm_source_dir/owl-$service_name.service"

  # 打包service_rpm_source_dir
  echo "[I] Tar service files"
  cd "$rpm_source_dir" && tar -cjf "owl-$service_name-$version.tar.bz2" "owl-$service_name-$version"

  # 开始构建rpm
  echo "[I] Start building rpm"
  cd "$rpmbuild_path/SPECS" && rpmbuild -ba "owl-$service_name-$version.spec"
  
  # 构建是否成功
  if [ $? -eq 0 ]; then
    echo "[I] Success build owl-$service_name-$version"
  else
    echo "[E] Failed build owl-$service_name-$version"
    exit 3
  fi
}

if [ $# -ne 2 ] ; then
  print_usage
fi

if [ $# -eq 2 ]; then
  build_go
  echo "[I] Build go finished"

  build_rpm
  echo "[I] Build rpm finished"

  exit 0
fi