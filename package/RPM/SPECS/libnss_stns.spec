%define _localbindir /usr/local/bin
%define _binaries_in_noarch_packages_terminate_build   0
Summary: SimpleTomlNameServiceLibrary
Name: libnss-stns
Group: SipmleTomlNameService
URL: https://github.com/pyama86/libnss_stns
Version: 0.0
Release: 3
License: MIT
Source0:   libnss_stns.conf
Packager:  libnss-stns
BuildRoot: %{_tmppath}/libnss-stns-%{version}-%{release}-root

%description
SimpleTomlNameService Client

%install
rm -rf %{buildroot}
install -d -m 755 %{buildroot}/%{_localbindir}
install    -m 655 %{_builddir}/stns-key-wrapper %{buildroot}/%{_localbindir}
install    -m 655 %{_builddir}/stns-query-wrapper %{buildroot}/%{_localbindir}

install -d -m 755 %{buildroot}/usr/%{_lib}
install    -m 655 %{_builddir}/libnss_stns.so %{buildroot}/usr/%{_lib}

install -d -m 755 %{buildroot}/%{_localstatedir}/log/

install -d -m 755 %{buildroot}/%{_sysconfdir}/stns/
install    -m 644 %{_sourcedir}/libnss_stns.conf %{buildroot}/%{_sysconfdir}/stns/libnss_stns.conf

%clean
rm -rf %{_builddir}
rm -rf %{buildroot}

%files
%defattr(-,root,root)
%{_localbindir}/stns-key-wrapper
%{_localbindir}/stns-query-wrapper
/usr/%{_lib}/libnss_stns.so
%config(noreplace) %{_sysconfdir}/stns/libnss_stns.conf

%post
ln -fs /usr/%{_lib}/libnss_stns.so /usr/%{_lib}/libnss_stns.so.2
