Summary:    build owl_kubecollect
Name:       owl-kubecollect
Version:    5.1.0
Release:    1%{?dist}
License:    Apache2.0
Group:      System Management
URL:        https://github.com/TalkingData/owl
%description

%pre
if [ "$1" = "2" ]; then
    # Perform whatever maintenance must occur before the upgrade begins
    systemctl stop owl-kubecollect
    cp ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/kubecollect.conf ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/kubecollect.conf.rpmsave
fi

%build
. /root/.bashrc
commitid=$(git -C $GOPATH/src/owl/ rev-parse --short HEAD)
go install -ldflags "-X main.Version=%{version} -X main.BuildTime=`date +%Y/%m/%d-%T` -X main.CommitID=${commitid}" owl/kubecollect

%install

#安装二进制
mkdir -p ${RPM_BUILD_ROOT}/usr/local/%{name}
cp -f $GOPATH/bin/kubecollect  ${RPM_BUILD_ROOT}/usr/local/%{name}/kubecollect

#服务管理文件
mkdir -p ${RPM_BUILD_ROOT}/usr/lib/systemd/system/
cp -f $GOPATH/src/owl/scripts/systemd/owl-kubecollect.service ${RPM_BUILD_ROOT}/usr/lib/systemd/system/owl-kubecollect.service

#服务配置文件
mkdir -p ${RPM_BUILD_ROOT}/usr/local/owl-kubecollect/conf
cp -f $GOPATH/src/owl/conf/kubecollect.conf  ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/kubecollect.conf

%post
if [ "$1" = "2" ];then
    cp ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/kubecollect.conf.rpmsave ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/kubecollect.conf
fi
systemctl daemon-reload
systemctl enable owl-kubecollect

%preun
systemctl stop owl-kubecollect

%postun
if [ "$1" = "0" ]; then
    # uninstall
    rm -rf ${RPM_BUILD_ROOT}/usr/local/%{name}
    systemctl daemon-reload
fi

%files
%defattr(-,root,root)
/usr/local/owl-kubecollect/kubecollect
/usr/lib/systemd/system/owl-kubecollect.service
/usr/local/owl-kubecollect/conf/kubecollect.conf
