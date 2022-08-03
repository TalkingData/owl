Summary:    build owl_netcollect
Name:       owl-netcollect
Version:    5.1.0
Release:    1%{?dist}
License:    Apache2.0
Group:      System Management
URL:        https://github.com/TalkingData/owl
Requires:   net-snmp-utils
%description

%pre
if [ "$1" = "2" ]; then
    # Perform whatever maintenance must occur before the upgrade begins
    systemctl stop owl-netcollect
    cp ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/netcollect.conf ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/netcollect.conf.rpmsave
fi

%build
. /root/.bashrc
commitid=$(git -C $GOPATH/src/owl/ rev-parse --short HEAD)
go install -ldflags "-X main.Version=%{version} -X main.BuildTime=`date +%Y/%m/%d-%T` -X main.CommitID=${commitid}" owl/netcollect

%install

#安装二进制
mkdir -p ${RPM_BUILD_ROOT}/usr/local/%{name}
cp -f $GOPATH/bin/netcollect  ${RPM_BUILD_ROOT}/usr/local/%{name}/netcollect

#服务管理文件
mkdir -p ${RPM_BUILD_ROOT}/usr/lib/systemd/system/
cp -f $GOPATH/src/owl/scripts/systemd/owl-netcollect.service ${RPM_BUILD_ROOT}/usr/lib/systemd/system/owl-netcollect.service

#服务配置文件
mkdir -p ${RPM_BUILD_ROOT}/usr/local/owl-netcollect/conf
cp -f $GOPATH/src/owl/conf/netcollect.conf  ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/netcollect.conf

%post
if [ "$1" = "2" ];then
    cp ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/netcollect.conf.rpmsave ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/netcollect.conf
fi
systemctl daemon-reload
systemctl enable owl-netcollect

%preun
systemctl stop owl-netcollect

%postun
if [ "$1" = "0" ]; then
    # uninstall
    rm -rf ${RPM_BUILD_ROOT}/usr/local/%{name}
    systemctl daemon-reload
fi

%files
%defattr(-,root,root)
/usr/local/owl-netcollect/netcollect
/usr/lib/systemd/system/owl-netcollect.service
/usr/local/owl-netcollect/conf/netcollect.conf
