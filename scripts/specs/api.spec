Summary:    build owl_api
Name:       owl-api
Version:    5.1.0
Release:    1%{?dist}
License:    Apache2.0
Group:      System Management
URL:        https://github.com/TalkingData/owl
Requires:   openssh,openssl
%description

%pre
if [ "$1" = "2" ]; then
    # Perform whatever maintenance must occur before the upgrade begins
    systemctl stop owl-api
    cp ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/api.conf ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/api.conf.rpmsave
fi

%build
. /root/.bashrc
commitid=$(git -C $GOPATH/src/owl/ rev-parse --short HEAD)
go install -ldflags "-X main.Version=%{version} -X main.BuildTime=`date +%Y/%m/%d-%T` -X main.CommitID=${commitid}" owl/api

%install

#安装二进制
mkdir -p ${RPM_BUILD_ROOT}/usr/local/%{name}
cp -f $GOPATH/bin/api  ${RPM_BUILD_ROOT}/usr/local/%{name}/api

#服务管理文件
mkdir -p ${RPM_BUILD_ROOT}/usr/lib/systemd/system/
cp -f $GOPATH/src/owl/scripts/systemd/owl-api.service ${RPM_BUILD_ROOT}/usr/lib/systemd/system/owl-api.service

#服务配置文件
mkdir -p ${RPM_BUILD_ROOT}/usr/local/owl-api/conf
cp -f $GOPATH/src/owl/conf/api.conf  ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/api.conf

mkdir -p  ${RPM_BUILD_ROOT}/usr/local/owl-api/certs
cd ${RPM_BUILD_ROOT}/usr/local/owl-api/certs
openssl genrsa -out rsa_private_key.pem 1024
openssl pkcs8 -topk8 -inform PEM -in rsa_private_key.pem -outform PEM -nocrypt -out owl-api.key
openssl rsa -in owl-api.key -pubout -out owl-api.key.pub



%post
if [ "$1" = "2" ];then
    cp ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/api.conf.rpmsave ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/api.conf
fi
systemctl daemon-reload
systemctl enable owl-api

%preun
systemctl stop owl-api

%postun
if [ "$1" = "0" ]; then
    # uninstall
    rm -rf ${RPM_BUILD_ROOT}/usr/local/%{name}
    systemctl daemon-reload
fi

%files
%defattr(-,root,root)
/usr/local/owl-api/api
/usr/lib/systemd/system/owl-api.service
/usr/local/owl-api/conf/api.conf
/usr/local/owl-api/certs/*
