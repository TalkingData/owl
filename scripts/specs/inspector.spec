Summary:    build owl_inspector
Name:       owl-inspector
Version:    5.1.0
Release:    1%{?dist}
License:    Apache2.0
Group:      System Management
URL:        https://github.com/TalkingData/owl
%description

%pre
if [ "$1" = "2" ]; then
    # Perform whatever maintenance must occur before the upgrade begins
    systemctl stop owl-inspector
    cp ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/inspector.conf ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/inspector.conf.rpmsave
fi

%build
. /root/.bashrc
commitid=$(git -C $GOPATH/src/owl/ rev-parse --short HEAD)
go install -ldflags "-X main.Version=%{version} -X main.BuildTime=`date +%Y/%m/%d-%T` -X main.CommitID=${commitid}" owl/inspector

%install

#安装二进制
mkdir -p ${RPM_BUILD_ROOT}/usr/local/%{name}
cp -f $GOPATH/bin/inspector  ${RPM_BUILD_ROOT}/usr/local/%{name}/inspector

#服务管理文件
mkdir -p ${RPM_BUILD_ROOT}/usr/lib/systemd/system/
cp -f $GOPATH/src/owl/scripts/systemd/owl-inspector.service ${RPM_BUILD_ROOT}/usr/lib/systemd/system/owl-inspector.service

#服务配置文件
mkdir -p ${RPM_BUILD_ROOT}/usr/local/owl-inspector/conf
cp -f $GOPATH/src/owl/conf/inspector.conf  ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/inspector.conf

mkdir -p  ${RPM_BUILD_ROOT}/usr/local/owl-inspector/certs
cd ${RPM_BUILD_ROOT}/usr/local/owl-inspector/certs



%post
if [ "$1" = "2" ];then
    cp ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/inspector.conf.rpmsave ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/inspector.conf
fi
systemctl daemon-reload
systemctl enable owl-inspector

%preun
systemctl stop owl-inspector

%postun
if [ "$1" = "0" ]; then
    # uninstall
    rm -rf ${RPM_BUILD_ROOT}/usr/local/%{name}
    systemctl daemon-reload
fi

%files
%defattr(-,root,root)
/usr/local/owl-inspector/inspector
/usr/lib/systemd/system/owl-inspector.service
/usr/local/owl-inspector/conf/inspector.conf
