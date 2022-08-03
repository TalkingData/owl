Summary:    build owl_proxy
Name:       owl-proxy
Version:    5.1.0
Release:    1%{?dist}
License:    Apache2.0
Group:      System Management
URL:        https://github.com/TalkingData/owl
%description

%pre
if [ "$1" = "2" ]; then
    # Perform whatever maintenance must occur before the upgrade begins
    systemctl stop owl-proxy
    cp ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/proxy.conf ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/proxy.conf.rpmsave
fi

%build
. /root/.bashrc
commitid=$(git -C $GOPATH/src/owl/ rev-parse --short HEAD)
go install -ldflags "-X main.Version=%{version} -X main.BuildTime=`date +%Y/%m/%d-%T` -X main.CommitID=${commitid}" owl/proxy

%install

#安装二进制
mkdir -p ${RPM_BUILD_ROOT}/usr/local/%{name}
cp -f $GOPATH/bin/proxy  ${RPM_BUILD_ROOT}/usr/local/%{name}/proxy

#服务管理文件
mkdir -p ${RPM_BUILD_ROOT}/usr/lib/systemd/system/
cp -f $GOPATH/src/owl/scripts/systemd/owl-proxy.service ${RPM_BUILD_ROOT}/usr/lib/systemd/system/owl-proxy.service

#服务配置文件
mkdir -p ${RPM_BUILD_ROOT}/usr/local/owl-proxy/conf
cp -f $GOPATH/src/owl/conf/proxy.conf  ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/proxy.conf

%post
if [ "$1" = "2" ];then
    cp ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/proxy.conf.rpmsave ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/proxy.conf
fi
systemctl daemon-reload
systemctl enable owl-proxy

%preun
systemctl stop owl-proxy

%postun
if [ "$1" = "0" ]; then
    # uninstall
    rm -rf ${RPM_BUILD_ROOT}/usr/local/%{name}
    systemctl daemon-reload
fi

%files
%defattr(-,root,root)
/usr/local/owl-proxy/proxy
/usr/lib/systemd/system/owl-proxy.service
/usr/local/owl-proxy/conf/proxy.conf
