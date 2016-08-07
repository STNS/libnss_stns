%define _localbindir /usr/local/bin
%define _stnslibdir /usr/lib/stns
%define _binaries_in_noarch_packages_terminate_build   0
Summary: Simple Toml NameService Library
Name: libnss-stns
Group: SipmleTomlNameService
URL: https://github.com/STNS/libnss_stns
Version: 0.1
Release: 8
License: MIT
Source0:   libnss_stns.conf
Packager:  libnss-stns
BuildRoot: %{_tmppath}/libnss-stns-%{version}-%{release}-root

%package -n libpam-stns
Summary: Simple Toml NameService Pam Module
Group: SipmleTomlNameService

%description -n libpam-stns
SimpleTomlNameService Pam Module

%description
SimpleTomlNameService Nss Module

%install
rm -rf %{buildroot}
install -d -m 755 %{buildroot}/%{_localbindir}
install    -m 755 %{_builddir}/stns-query-wrapper %{buildroot}/%{_localbindir}

install -d -m 755 %{buildroot}/%{_stnslibdir}
install    -m 755 %{_builddir}/stns-key-wrapper %{buildroot}/%{_stnslibdir}

ln -fs %{_stnslibdir}/stns-key-wrapper %{buildroot}%{_localbindir}

install -d -m 755 %{buildroot}/usr/%{_lib}
install    -m 755 %{_builddir}/libnss_stns.so %{buildroot}/usr/%{_lib}

install -d -m 755 %{buildroot}/%{_lib}/security
install    -m 755 %{_builddir}/libpam_stns.so %{buildroot}/%{_lib}/security

install -d -m 755 %{buildroot}/%{_localstatedir}/log/

install -d -m 755 %{buildroot}/%{_sysconfdir}/stns/
install    -m 644 %{_sourcedir}/libnss_stns.conf %{buildroot}/%{_sysconfdir}/stns/libnss_stns.conf
ln -fs /usr/%{_lib}/libnss_stns.so %{buildroot}/usr/%{_lib}/libnss_stns.so.2

install -d -m 777 %{buildroot}/var/lib/libnss_stns

%clean
rm -rf %{_builddir}
rm -rf %{buildroot}

%files
%defattr(-,root,root)
%{_stnslibdir}/stns-key-wrapper
%{_localbindir}/stns-key-wrapper
%{_localbindir}/stns-query-wrapper
/usr/%{_lib}/libnss_stns.so
/usr/%{_lib}/libnss_stns.so.2
%config(noreplace) %{_sysconfdir}/stns/libnss_stns.conf

%dir %attr(0700, root, root) %{_stnslibdir}
%dir /var/lib/libnss_stns

%files -n libpam-stns
/%{_lib}/security/libpam_stns.so

%post
ln -fs /usr/%{_lib}/libnss_stns.so /usr/%{_lib}/libnss_stns.so.2
