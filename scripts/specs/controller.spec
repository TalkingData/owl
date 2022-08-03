Summary:    build owl_controller
Name:       owl-controller
Version:    5.1.0
Release:    1%{?dist}
License:    Apache2.0
Group:      System Management
URL:        https://github.com/TalkingData/owl
Requires:   python
%description

%pre
if [ "$1" = "2" ]; then
    # Perform whatever maintenance must occur before the upgrade begins
    systemctl stop owl-controller
    cp ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/controller.conf ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/controller.conf.rpmsave
fi

%build
. /root/.bashrc
commitid=$(git -C $GOPATH/src/owl/ rev-parse --short HEAD)
go install -ldflags "-X main.Version=%{version} -X main.BuildTime=`date +%Y/%m/%d-%T` -X main.CommitID=${commitid}" owl/controller

%install

#安装二进制
mkdir -p ${RPM_BUILD_ROOT}/usr/local/%{name}
cp -f $GOPATH/bin/controller  ${RPM_BUILD_ROOT}/usr/local/%{name}/controller

#服务管理文件
mkdir -p ${RPM_BUILD_ROOT}/usr/lib/systemd/system/
cp -f $GOPATH/src/owl/scripts/systemd/owl-controller.service ${RPM_BUILD_ROOT}/usr/lib/systemd/system/owl-controller.service

#服务配置文件
mkdir -p ${RPM_BUILD_ROOT}/usr/local/owl-controller/conf
cp -f $GOPATH/src/owl/conf/controller.conf  ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/controller.conf

#脚本文件
cp -rf $GOPATH/src/owl/controller/scripts  ${RPM_BUILD_ROOT}/usr/local/%{name}/


%post
if [ "$1" = "2" ];then
    cp ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/controller.conf.rpmsave ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/controller.conf
fi
systemctl daemon-reload
systemctl enable owl-controller

%preun
systemctl stop owl-controller

%postun
if [ "$1" = "0" ]; then
    # uninstall
    rm -rf ${RPM_BUILD_ROOT}/usr/local/%{name}
    systemctl daemon-reload
fi

%files
%defattr(-,root,root)
/usr/local/owl-controller/controller
/usr/lib/systemd/system/owl-controller.service
/usr/local/owl-controller/conf/controller.conf
/usr/local/owl-controller/scripts/*
