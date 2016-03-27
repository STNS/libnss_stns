%define _localbindir /usr/local/bin
%define _binaries_in_noarch_packages_terminate_build   0
Summary: SimpleTomlNameServiceLibrary
Name: lib-stns
Group: SipmleTomlNameService
URL: https://github.com/STNS/lib-stns
Version: 0.0
Release: 5
License: MIT
Source0:   lib_stns.conf
Packager:  lib-stns
BuildRoot: %{_tmppath}/lib-stns-%{version}-%{release}-root
Obsoletes: libnss-stns

%description
SimpleTomlNameService Client

%install
rm -rf %{buildroot}
install -d -m 755 %{buildroot}/%{_localbindir}
install    -m 655 %{_builddir}/stns-key-wrapper %{buildroot}/%{_localbindir}
install    -m 655 %{_builddir}/stns-query-wrapper %{buildroot}/%{_localbindir}

install -d -m 755 %{buildroot}/usr/%{_lib}
install    -m 655 %{_builddir}/libnss-stns.so %{buildroot}/usr/%{_lib}/libnss_stns.so

install -d -m 755 %{buildroot}/%{_lib}/security
install    -m 655 %{_builddir}/libpam-stns.so %{buildroot}/%{_lib}/security/libpam_stns.so

install -d -m 755 %{buildroot}/%{_localstatedir}/log/

install -d -m 755 %{buildroot}/%{_sysconfdir}/stns/
install    -m 644 %{_sourcedir}/lib_stns.conf %{buildroot}/%{_sysconfdir}/stns/lib_stns.conf
ln -fs /usr/%{_lib}/libnss_stns.so %{buildroot}/usr/%{_lib}/libnss_stns.so.2

%clean
rm -rf %{_builddir}
rm -rf %{buildroot}

%files
%defattr(-,root,root)
%{_localbindir}/stns-key-wrapper
%{_localbindir}/stns-query-wrapper
/usr/%{_lib}/libnss_stns.so
/usr/%{_lib}/libnss_stns.so.2
/%{_lib}/security/libpam_stns.so
%config(noreplace) %{_sysconfdir}/stns/lib_stns.conf

%post
ln -fs /usr/%{_lib}/libnss_stns.so /usr/%{_lib}/libnss_stns.so.2

if [ -e /etc/stns/libnss_stns.conf && ! -e /etc/stns/lib_stns.conf ]; then
  cp -p /etc/stns/libnss_stns.conf /etc/stns/libnss_stns.conf.back
  mv /etc/stns/libnss_stns.conf /etc/stns/lib_stns.conf
  echo "move config file libnss_stns.conf to lib_stns.conf"
fi

%preun
if [ $1 = 0 ]; then
  rm -rf /usr/%{_lib}/libnss_stns.so.2
fi
