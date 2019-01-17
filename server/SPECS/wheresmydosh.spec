Name: wheresmydosh 
Version: 0.0.1
Release: %{?_buildnumber}%{?dist}
Group: Applications/Tools
License: GPL
Packager: Kodjo Baah <kodjo_baah@hotmail.com>
Summary: Api used to share money between friends and family
#Source: wheresmydosh.tar.gz
BuildRoot: %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)
BuildArch: x86_64


Requires: systemd
Requires: lighttpd
Requires: lighttpd-fastcgi
Requires: lighttpd-mod_geoip
Requires: mysql-server
Requires: mysql


%description
This is used for processing file upload

%build
mkdir -p %{buildroot}/%{_bindir}/
mkdir -p %{buildroot}/etc/systemd/system/
mkdir -p %{buildroot}/etc/wheresmydosh


cp %{_topdir}/bin/wheresmydosh %{buildroot}/%{_bindir}
cp %{_topdir}/scripts/run_wheresmydosh.sh %{buildroot}/%{_bindir}
cp %{_topdir}/init/systemd/wheresmydosh.service %{buildroot}/etc/systemd/system
cp %{_topdir}/config/etc/wheresmydosh/wheresmydosh.conf %{buildroot}/etc/wheresmydosh


%files
%defattr(-,root,root)
%{_bindir}/wheresmydosh
%{_bindir}/run_wheresmydosh.sh
/etc/systemd/system/wheresmydosh.service
/etc/wheresmydosh/wheresmydosh.conf


%post
%systemd_post wheresmydosh.service

%preun
%systemd_preun wheresmydosh.service

%postun
%systemd_postun wheresmydosh.service
