Name:      owl-repeater
Version:   6.0.2
Release:   1%{?dist}
Summary:   TalkingData OWL Repeater

Group:      TalkingData OWL
License:    Apache2.0
URL:        https://github.com/TalkingData/owl
Packager:   lai.li <lai.li@tendcloud.com>
Source0:    %{name}-%{version}.tar.bz2
BuildRoot:  %_topdir/BUILDROOT
Prefix:     /usr/local/%{name}

Requires: tar,bzip2

%description
TalkingData OWL Repeater...

%prep
%setup -q

# 安装阶段
%install
mkdir -p ${RPM_BUILD_ROOT}/usr/local/%{name}/{bin,conf,logs}
cp -f %_topdir/BUILD/%{name}-%{version}/bin/owl-repeater ${RPM_BUILD_ROOT}/usr/local/%{name}/bin/owl-repeater
cp -f %_topdir/BUILD/%{name}-%{version}/conf/owl_repeater.conf ${RPM_BUILD_ROOT}/usr/local/%{name}/conf/owl_repeater.conf
mkdir -p ${RPM_BUILD_ROOT}/usr/lib/systemd/system/
cp -f %_topdir/BUILD/%{name}-%{version}/owl-repeater.service ${RPM_BUILD_ROOT}/usr/lib/systemd/system/owl-repeater.service

# 安装前执行的脚本，语法和shell脚本的语法相同
%pre
if [ $1 == 1 ];then
echo -e "\e[1;31mInstall to /usr/local/owl-repeater/\e[0m"
fi
if [ $1 == 2 ];then
echo -e "\e[1;31mReplace /usr/local/owl-repeater/\e[0m"
fi

# 安装后执行的脚本
%post
if [ $1 == 1 ];then
systemctl daemon-reload
fi

# 卸载前执行的脚本
%preun
systemctl stop owl-repeater
systemctl disable owl-repeater
/usr/bin/cp /usr/local/owl-repeater/conf/owl_repeater.conf /usr/local/owl-repeater/conf/owl_repeater.conf.rpmsave

# 卸载完成后执行的脚本
%postun
if [ $1 == 0 ];then
echo -e "\e[1;32m%{name} Uninstalled.\e[0m"
rm -rf /usr/local/owl-repeater/bin/owl-repeater
rm -rf /usr/local/owl-repeater/conf/owl_repeater.conf
rm -rf /usr/lib/systemd/system/owl-repeater.service
fi
if [ $1 == 1 ];then
echo "Removed."
fi

# 清理阶段，在制作完成后删除安装的内容
%clean
rm -rf %{buildroot}

#指定要包含的文件
%files

#设置默认权限，如果没有指定，则继承默认的权限
%defattr(-,root,root,0644)
/usr/local/owl-repeater/bin/owl-repeater
/usr/local/owl-repeater/conf/owl_repeater.conf
/usr/lib/systemd/system/owl-repeater.service
