%define _localbindir /usr/local/bin
%define _binaries_in_noarch_packages_terminate_build   0
Summary: SimpleTomlNameServiceLibrary
Name: libnss_stns
Group: SipmleTomlNameService
URL: https://github.com/pyama86/libnss_stns
Version: 0.1
Release: 1
License: MIT
Source0:   %{name}.conf
Packager:  libnss_stns
BuildArch: noarch
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
SimpleTomlNameService Client

%install
rm -rf %{buildroot}
install -d -m 755 %{buildroot}/%{_localbindir}
install    -m 655 %{_builddir}/ssh_stns_wrapper %{buildroot}/%{_localbindir}

install -d -m 755 %{buildroot}/usr/%{_lib}
install    -m 655 %{_builddir}/%{name}.so %{buildroot}/usr/%{_lib}

install -d -m 755 %{buildroot}/%{_localstatedir}/log/

install -d -m 755 %{buildroot}/%{_sysconfdir}/stns/
install    -m 644 %{_sourcedir}/%{name}.conf %{buildroot}/%{_sysconfdir}/stns/%{name}.conf

%clean
rm -rf %{_builddir}
rm -rf %{buildroot}

%files
%defattr(-,root,root)
%{_localbindir}/ssh_stns_wrapper
/usr/%{_lib}/%{name}.so
%config(noreplace) %{_sysconfdir}/stns/%{name}.conf
